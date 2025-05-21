package firebase

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

func (f *Firebase) VerifyIdToken(ctx context.Context, idToken string) (string, error) {
	token, err := f.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}

func (f *Firebase) DeleteAccount(ctx context.Context, firebaseUid string) error {
	return f.auth.DeleteUser(ctx, firebaseUid)
}

func (f *Firebase) SendNotification(ctx context.Context, token string, title, body string, data map[string]string) (string, error) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Token: token,
	}
	return f.messaging.Send(ctx, message)
}

func (f *Firebase) SendMulticastNotification(ctx context.Context, tokens []string, title, body string, data map[string]string) (*messaging.BatchResponse, error) {
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
	}
	return f.messaging.SendEachForMulticast(ctx, message)
}
