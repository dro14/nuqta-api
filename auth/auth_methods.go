package auth

import "context"

func (a *Auth) VerifyIdToken(ctx context.Context, idToken string) (string, error) {
	token, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}
