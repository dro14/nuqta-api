package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Firebase struct {
	auth      *auth.Client
	messaging *messaging.Client
}

func New() *Firebase {
	ctx := context.Background()
	opt := option.WithCredentialsFile("service_account_key.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal("can't initialize firebase app: ", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.Fatal("can't initialize firebase auth: ", err)
	}

	messaging, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal("can't initialize firebase messaging: ", err)
	}

	return &Firebase{
		auth:      auth,
		messaging: messaging,
	}
}
