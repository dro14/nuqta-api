package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dro14/nuqta-service/e"
)

func (d *Data) GetReferrer(ctx context.Context, ip, osVersion string) (string, error) {
	osVersion = strings.TrimSuffix(strings.TrimSuffix(osVersion, ".0"), ".0")
	key := fmt.Sprintf("referrer:%s:%s", ip, osVersion)
	referrer, err := d.cacheGet(key)
	if errors.Is(err, e.ErrNotFound) {
		key = fmt.Sprintf("referrer:%s", ip)
		referrer, err = d.cacheGet(key)
		if errors.Is(err, e.ErrNotFound) {
			return "", nil
		} else if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	userUid := string(referrer)
	var registered int64
	err = d.dbQueryRow(ctx, "SELECT registered FROM users WHERE id = $1", []any{userUid}, &registered)
	if errors.Is(err, e.ErrNotFound) || registered == 0 {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return userUid, nil
}

func (d *Data) SetReferrer(ctx context.Context, ip, osVersion, referrer string) error {
	osVersion = strings.TrimSuffix(strings.TrimSuffix(osVersion, ".0"), ".0")
	key := fmt.Sprintf("referrer:%s", ip)
	if osVersion != "" {
		key += fmt.Sprintf(":%s", osVersion)
	}
	return d.cacheSet(key, []byte(referrer), 24*time.Hour)
}
