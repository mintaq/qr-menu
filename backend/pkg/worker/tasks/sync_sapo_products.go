package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeSyncSapoProducts = "sync:sapo:products"

type SyncSapoProductsPayload struct {
	SinceId   int
	Limit     int
	UserAppId uint64
}

func NewSyncSapoProductsRecursiveTask(sinceId int, limit int, userAppId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncSapoProductsPayload{SinceId: sinceId, Limit: limit, UserAppId: userAppId})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncSapoProducts, payload), nil
}

func HandleSyncSapoProductsRecursiveTask(ctx context.Context, t *asynq.Task) error {

	return nil
}
