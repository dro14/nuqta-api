package auth

import "context"

func (a *Auth) VerifyIdToken(idToken string) (string, error) {
	token, err := a.client.VerifyIDToken(context.Background(), idToken)
	return token.UID, err
}
