package data

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/dro14/nuqta-service/models"
	"github.com/dro14/nuqta-service/utils/e"
	"github.com/lib/pq"
)

const retryAttempts = 5

func getSearchVector(name, username string) string {
	searchVector := ""
	if name != "" {
		searchVector += strings.TrimSpace(name) + " "
	}
	if username != "" {
		searchVector += "@" + strings.TrimSpace(username)
	}
	return searchVector
}

func (d *Data) graphGet(ctx context.Context, query string, vars map[string]string) ([]byte, error) {
	var lastErr error
	for range retryAttempts {
		txn := d.graph.NewReadOnlyTxn().BestEffort()
		resp, err := txn.QueryWithVars(ctx, query, vars)
		if err == nil {
			return resp.Json, nil
		}
		lastErr = err
		log.Printf("can't get: %s%s", err, query)
	}
	log.Printf("failed to get after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Data) graphSet(ctx context.Context, object any) (*api.Response, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	mutation := &api.Mutation{
		SetJson:   bytes,
		CommitNow: true,
	}
	var lastErr error
	for range retryAttempts {
		txn := d.graph.NewTxn()
		resp, err := txn.Mutate(ctx, mutation)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		log.Printf("can't set: %s\n%s", err, bytes)
	}
	log.Printf("failed to set after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Data) graphDelete(ctx context.Context, object any) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	mutation := &api.Mutation{
		DeleteJson: bytes,
		CommitNow:  true,
	}
	var lastErr error
	for range retryAttempts {
		txn := d.graph.NewTxn()
		_, err := txn.Mutate(ctx, mutation)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("can't delete: %s\n%s", err, bytes)
	}
	log.Printf("failed to delete after %d attempts", retryAttempts)
	return lastErr
}

func (d *Data) dbExec(ctx context.Context, query string, args ...any) error {
	var lastErr error
	for range retryAttempts {
		_, err := d.db.ExecContext(ctx, query, args...)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("can't exec: %s\n%s", err, query)
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrNotFound
		}
	}
	log.Printf("failed to exec after %d attempts", retryAttempts)
	return lastErr
}

func (d *Data) dbQueryRow(ctx context.Context, query string, args []any, dest ...any) error {
	var lastErr error
	for range retryAttempts {
		err := d.db.QueryRowContext(ctx, query, args...).Scan(dest...)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("can't query row: %s\n%s", err, query)
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrNotFound
		}
	}
	log.Printf("failed to query row after %d attempts", retryAttempts)
	return lastErr
}

func (d *Data) dbQuery(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var lastErr error
	for range retryAttempts {
		rows, err := d.db.QueryContext(ctx, query, args...)
		if err == nil {
			return rows, nil
		}
		lastErr = err
		log.Printf("can't query: %s\n%s", err, query)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
	}
	log.Printf("failed to query after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Data) cacheGet(key string) ([]byte, error) {
	item, err := d.cache.Get(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil, e.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (d *Data) cacheSet(key string, value []byte, ttl time.Duration) error {
	return d.cache.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	})
}

func (d *Data) cacheDelete(key string) error {
	err := d.cache.Delete(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil
	}
	return err
}

func decodeMessages(rows *sql.Rows, type_ string) []*models.Message {
	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := &models.Message{}
		var nullInReplyTo sql.NullInt64
		var nullText sql.NullString
		var nullViewed sql.NullInt64
		var nullLiked sql.NullInt64
		var nullEdited sql.NullInt64
		var nullDeleted sql.NullInt64
		dest := []any{&message.Id, &message.Timestamp, &message.ChatUid, &message.AuthorUid, &nullInReplyTo, &nullText, pq.Array(&message.Images)}
		if type_ == "private" {
			dest = append(dest, &nullViewed, &nullLiked, &nullEdited, &nullDeleted, &message.RecipientUid)
		}
		err := rows.Scan(dest...)
		if err != nil {
			log.Printf("can't scan %s message: %s", type_, err)
			continue
		}
		if nullInReplyTo.Valid {
			message.InReplyTo = nullInReplyTo.Int64
		}
		if nullText.Valid {
			message.Text = nullText.String
		}
		if nullViewed.Valid {
			message.Viewed = nullViewed.Int64
		}
		if nullLiked.Valid {
			message.Liked = nullLiked.Int64
		}
		if nullEdited.Valid {
			message.Edited = nullEdited.Int64
		}
		if nullDeleted.Valid {
			message.Deleted = nullDeleted.Int64
		}
		messages = append(messages, message)
	}
	return messages
}

func decodeIds(rows *sql.Rows) []int64 {
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("can't scan id: %v", err)
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func scoreUser(user *models.User) int {
	return 2*user.Followers +
		1*user.Invites -
		2*user.Blockers
}

func scorePost(post *models.Post) int {
	return 20*post.Replies +
		15*post.Reposts +
		10*post.Likes +
		5*post.Clicks +
		1*post.Views
}

func (d *Data) deleteImages(ctx context.Context, images []string) error {
	if len(images) == 0 {
		return nil
	}
	for i, image := range images {
		_, images[i], _ = strings.Cut(image, "/images/")
	}
	request := map[string][]string{"filenames": images}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		return err
	}
	url := "https://images.nuqtam.uz/hammasi"
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, &buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(d.username, d.password)

	var client http.Client
	client.Timeout = 10 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bts))
	}
	return nil
}
