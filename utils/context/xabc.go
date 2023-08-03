package context

import (
	"context"
)

type xAbcKey struct{}

func SetXAbc(ctx context.Context, xAbc string) context.Context {
	return context.WithValue(ctx, xAbcKey{}, xAbc)
}

func GetXAbc(ctx context.Context) *string {
	xAbc := ctx.Value(xAbcKey{})
	if xAbc == nil {
		return nil
	}

	k := xAbc.(string)
	return &k
}
