package mongo

import "context"

func id(ctx context.Context) string {
	return ctx.Value("id").(string)
}
