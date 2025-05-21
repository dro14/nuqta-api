package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/dro14/nuqta-service/utils/e"
	"github.com/lib/pq"
)

func (d *Data) CreateProfile(ctx context.Context, user *models.User) error {
	suffixMinLength := 0
	username, _ := generateUsername(user.Email, suffixMinLength)
	for {
		userUid, err := d.GetUidByUsername(ctx, username)
		if userUid != "" {
			suffixMinLength++
			username, err = generateUsername(user.Email, suffixMinLength)
			if err != nil {
				return err
			}
		} else if errors.Is(err, e.ErrNotFound) {
			user.Username = username
			break
		} else {
			return err
		}
	}
	user.Registered = time.Now().UnixMilli()

	object := map[string]any{
		"dgraph.type": "user",
		"uid":         "_:user",
		"registered":  user.Registered,
	}

	if user.InvitedBy != nil {
		object["invited_by"] = map[string]string{
			"uid": user.InvitedBy.Uid,
		}
	}

	assigned, err := d.graphSet(ctx, object)
	if err != nil {
		return err
	}
	user.Uid = assigned.Uids["user"]

	err = d.dbExec(ctx,
		"INSERT INTO users (id, firebase_uid, email, registered, name, username, search_vector) VALUES ($1, $2, $3, $4, $5, $6, to_tsvector('simple', $7))",
		user.Uid, user.FirebaseUid, user.Email, user.Registered, user.Name, user.Username, getSearchVector(user.Name, user.Username),
	)
	if err != nil {
		object = map[string]any{"uid": user.Uid}
		d.graphDelete(ctx, object)
		return err
	}

	return nil
}

func (d *Data) GetUidByFirebaseUid(ctx context.Context, firebaseUid string) (string, error) {
	bytes, err := d.cacheGet("uid:" + firebaseUid)
	if errors.Is(err, e.ErrNotFound) {
		var uid string
		err := d.dbQueryRow(ctx,
			"SELECT id FROM users WHERE firebase_uid = $1",
			[]any{firebaseUid},
			&uid,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrNotFound
		} else if err != nil {
			return "", err
		}
		return uid, d.cacheSet("uid:"+firebaseUid, []byte(uid), 0)
	}
	return string(bytes), err
}

func (d *Data) GetProfile(ctx context.Context, firebaseUid string) (*models.User, error) {
	uid, err := d.GetUidByFirebaseUid(ctx, firebaseUid)
	if err != nil {
		return nil, err
	}

	user, err := d.GetUser(ctx, uid, uid)
	if err != nil {
		return nil, err
	}

	user.FirebaseUid = firebaseUid
	return user, nil
}

func (d *Data) EditProfile(ctx context.Context, user *models.User) error {
	return d.dbExec(ctx,
		"UPDATE users SET name = $1, username = $2, location = $3, birthday = $4, color = $5, bio = $6, banner = $7, avatars = $8, thumbnails = $9, search_vector = to_tsvector('simple', $10) WHERE id = $11",
		user.Name, user.Username, user.Location, user.Birthday, user.Color, user.Bio, user.Banner, pq.Array(user.Avatars), pq.Array(user.Thumbnails), getSearchVector(user.Name, user.Username), user.Uid,
	)
}

func (d *Data) DeleteProfile(ctx context.Context, uid, firebaseUid string) error {
	postUids, err := d.GetUserPosts(ctx, "all", uid, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	for _, postUid := range postUids {
		post, err := d.GetPost(ctx, uid, postUid)
		if err != nil {
			log.Printf("can't get post: %s", err)
			continue
		}
		err = d.DeletePost(ctx, uid, postUid, post.Images)
		if err != nil {
			log.Printf("can't delete post: %s", err)
		}
	}

	user, err := d.GetUser(ctx, uid, uid)
	if err != nil {
		return err
	}

	err = d.deleteImages(ctx, user.Avatars)
	if err != nil {
		log.Printf("can't delete avatars: %s", err)
	}

	err = d.deleteImages(ctx, user.Thumbnails)
	if err != nil {
		log.Printf("can't delete thumbnails: %s", err)
	}

	err = d.cacheDelete("uid:" + firebaseUid)
	if err != nil {
		log.Printf("can't delete uid from cache: %s", err)
	}

	object := map[string]any{"uid": uid}
	err = d.graphDelete(ctx, object)
	if err != nil {
		return err
	}
	return d.dbExec(ctx, "DELETE FROM users WHERE id = $1", uid)
}
