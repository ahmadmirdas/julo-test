package activity

import (
	"context"

	"github.com/google/uuid"
)

type key int

const (
	ActivityID key = iota
	Action
	ActorID
	ActorIP
	Actor
	CakeID
)

func NewContext(action string) context.Context {
	activityId := uuid.New().String()
	ctx := context.WithValue(context.Background(), ActivityID, activityId)
	return context.WithValue(ctx, Action, action)
}

func GetActivityID(ctx context.Context) (string, bool) {
	return getStringValueFromContext(ctx, ActivityID)
}

func GetAction(ctx context.Context) (string, bool) {
	return getStringValueFromContext(ctx, Action)
}

func WithCakeID(ctx context.Context, cakeID int) context.Context {
	return context.WithValue(ctx, CakeID, cakeID)
}

func GetFields(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	if id, ok := GetActivityID(ctx); ok {
		fields["activity_id"] = id
	}
	if action, ok := GetAction(ctx); ok {
		fields["action"] = action
	}
	return fields
}

func getStringValueFromContext(ctx context.Context, key key) (string, bool) {
	value, ok := ctx.Value(key).(string)
	return value, ok
}

func getIntValueFromContext(ctx context.Context, key key) (int, bool) {
	value, ok := ctx.Value(key).(int)
	return value, ok
}
