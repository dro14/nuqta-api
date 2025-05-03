package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/lib/pq"
)

func (d *Data) GetUser(ctx context.Context, uid, userUid string) (*models.User, error) {
	vars := map[string]string{
		"$user_uid": userUid,
	}
	bytes, err := d.graphGet(ctx, userQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.User
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	user := response["users"][0]
	user.Uid = userUid

	var nullLocation sql.NullString
	var nullBirthday sql.NullInt64
	var nullColor sql.NullString
	var nullBio sql.NullString
	var nullBanner sql.NullString

	err = d.dbQueryRow(ctx,
		"SELECT registered, name, username, location, birthday, color, bio, banner, avatars, thumbnails FROM users WHERE id = $1",
		[]any{userUid},
		&user.Registered, &user.Name, &user.Username, &nullLocation, &nullBirthday, &nullColor, &nullBio, &nullBanner, pq.Array(&user.Avatars), pq.Array(&user.Thumbnails),
	)
	if err != nil {
		return nil, err
	}

	if nullLocation.Valid {
		user.Location = nullLocation.String
	}
	if nullBirthday.Valid {
		user.Birthday = nullBirthday.Int64
	}
	if nullColor.Valid {
		user.Color = nullColor.String
	}
	if nullBio.Valid {
		user.Bio = nullBio.String
	}
	if nullBanner.Valid {
		user.Banner = nullBanner.String
	}

	if uid == userUid {
		return user, nil
	}

	vars = map[string]string{
		"$uid":      uid,
		"$user_uid": userUid,
	}
	bytes, err = d.graphGet(ctx, userEdgesQuery, vars)
	if err != nil {
		return nil, err
	}

	var edges map[string][]map[string][]map[string]string
	err = json.Unmarshal(bytes, &edges)
	if err != nil {
		return nil, err
	}

	if len(edges["users"]) > 0 {
		user_ := edges["users"][0]
		user.IsFollowed = len(user_["is_followed"]) > 0
		user.IsFollowing = len(user_["is_following"]) > 0
		if len(user_["chats"]) == 1 {
			user.ChatUid = user_["chats"][0]["uid"]
		}
	}

	return user, nil
}

func (d *Data) GetUidByUsername(ctx context.Context, username string) (string, error) {
	var userUid string
	err := d.dbQueryRow(ctx,
		"SELECT id FROM users WHERE LOWER(username) = $1",
		[]any{strings.ToLower(username)},
		&userUid,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return "", e.ErrNotFound
	} else if err != nil {
		return "", err
	}
	return userUid, nil
}

func (d *Data) SearchUser(ctx context.Context, query string, offset int64) ([]string, error) {
	query = strings.Join(strings.Fields(query), ":* & ") + ":*"
	rows, err := d.dbQuery(ctx, `
		SELECT id, ts_rank(search_vector, query) AS rank
		FROM users, to_tsquery('simple', $1) query
		WHERE search_vector @@ query
		ORDER BY rank DESC
		LIMIT 20
		OFFSET $2;`,
		query, offset,
	)
	if err != nil {
		return nil, err
	}

	var userUids []string
	for rows.Next() {
		var userUid string
		var rank float64
		err = rows.Scan(&userUid, &rank)
		if err != nil {
			log.Printf("can't scan: %s", err)
			continue
		}
		userUids = append(userUids, userUid)
	}

	return userUids, nil
}

func (d *Data) GetUserFollows(ctx context.Context, userUid string, offset int64, reverse bool) ([]string, error) {
	vars := map[string]string{
		"$user_uid": userUid,
		"$offset":   strconv.FormatInt(offset, 10),
	}
	var edge string
	if reverse {
		edge = "~follow"
	} else {
		edge = "follow"
	}
	query := fmt.Sprintf(userFollowsQuery, edge)
	bytes, err := d.graphGet(ctx, query, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.User
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var userUids []string
	for _, user := range response["users"] {
		for _, follower := range user[edge] {
			userUids = append(userUids, follower.Uid)
		}
	}

	return userUids, nil
}

func (d *Data) GetUserInvites(ctx context.Context, uid string, offset int64) ([]string, error) {
	vars := map[string]string{
		"$uid":    uid,
		"$offset": strconv.FormatInt(offset, 10),
	}
	bytes, err := d.graphGet(ctx, userInvitesQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.User
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var userUids []string
	for _, user := range response["users"] {
		for _, invited := range user["invited"] {
			userUids = append(userUids, invited.Uid)
		}
	}

	return userUids, nil
}
