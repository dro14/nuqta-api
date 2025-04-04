package yordamchi

import (
	"context"
	"log"
	"time"
)

func sleep(retryDelay *time.Duration) {
	if *retryDelay > 0 {
		log.Print("retrying request after ", *retryDelay)
		time.Sleep(*retryDelay)
		*retryDelay *= 2
	}
}

func id(ctx context.Context) string {
	return ctx.Value("firebase_uid").(string)
}
