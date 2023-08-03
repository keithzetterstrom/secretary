package context

import "context"

type requestIDKey struct{}

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, reqID)
}

func GetRequestID(ctx context.Context) *string {
	requestID := ctx.Value(requestIDKey{})
	if requestID == nil {
		return nil
	}

	id := requestID.(string)
	return &id
}
