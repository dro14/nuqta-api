package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
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
	user.Registered = time.Now().Unix()

	object := map[string]any{
		"dgraph.type": "User",
		"uid":         "_:user",
	}

	assigned, err := d.graphSet(ctx, object)
	if err != nil {
		return err
	}
	user.Uid = assigned.Uids["user"]

	err = d.dbExec(ctx,
		"INSERT INTO users (id, firebase_uid, email, registered, username, name) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Uid, user.FirebaseUid, user.Email, user.Registered, user.Username, user.Name,
	)
	if err != nil {
		object = map[string]any{"uid": user.Uid}
		d.graphDelete(ctx, object)
		return err
	}

	return nil
}

func (d *Data) GetProfile(ctx context.Context, firebaseUid string) (*models.User, error) {
	var uid string
	err := d.dbQueryRow(ctx,
		"SELECT id FROM users WHERE firebase_uid = $1",
		firebaseUid,
		&uid,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, e.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	user, err := d.GetUser(ctx, uid, uid)
	if err != nil {
		return nil, err
	}

	user.FirebaseUid = firebaseUid
	return user, nil
}

func (d *Data) UpdateProfile(ctx context.Context, user *models.User) error {
	return d.dbExec(ctx,
		"UPDATE users SET username = $1, name = $2, location = $3, birthday = $4, color = $5, bio = $6, banner = $7, avatars = $8, thumbnails = $9 WHERE id = $10",
		user.Username, user.Name, user.Location, user.Birthday, user.Color, user.Bio, user.Banner, pq.Array(user.Avatars), pq.Array(user.Thumbnails), user.Uid,
	)
}
