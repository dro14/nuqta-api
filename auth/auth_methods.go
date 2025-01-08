package auth

import "context"

func (a *Auth) VerifyIdToken(idToken string) (string, error) {
	token, err := a.client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}
