package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Auth struct {
	client *auth.Client
}

func New() *Auth {
	ctx := context.Background()
	opt := option.WithCredentialsFile("service_account_key.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal("can't initialize firebase app: ", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatal("can't initialize firebase auth: ", err)
	}

	return &Auth{client: client}
}
