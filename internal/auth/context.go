package auth

import "context"

// An unexported zero-size struct type means no other package can construct this key, so nothing outside
// this package can collide with (or overwrite) the value we store under it
type userIDContextKey struct{}

func ContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey{}).(int64)
	return userID, ok
}
