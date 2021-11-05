package blogcontext

import (
	"context"
	"errors"
	"fmt"
)

type blogContextKey string

var (
	ErrValueNotPresent = errors.New("context doesn't have the requested value")
)

func AddRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) (string, error) {
	value := ctx.Value(requestIDKey)
	if value == nil {
		return "", ErrValueNotPresent
	}

	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("failed typecasting %T to string", value)
	}

	return s, nil
}
