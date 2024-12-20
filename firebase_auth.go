package main

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(credentialsFile string) (*FirebaseAuth, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseAuth{client: client}, nil
}

func (fa *FirebaseAuth) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		idToken := strings.TrimPrefix(header, "Bearer ")

		token, err := fa.client.VerifyIDToken(ctx, idToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
