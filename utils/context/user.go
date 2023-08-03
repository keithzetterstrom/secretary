package context

import "context"

type userIDKey struct{}
type userNameKey struct{}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) *int64 {
	userID := ctx.Value(userIDKey{})
	if userID == nil {
		return nil
	}

	id := userID.(int64)
	return &id
}

func SetUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userNameKey{}, userName)
}

func GetUserName(ctx context.Context) *string {
	userName := ctx.Value(userNameKey{})
	if userName == nil {
		return nil
	}

	name := userName.(string)
	return &name
}
