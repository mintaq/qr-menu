package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/sapo"
)

const (
	TypeSyncSapoCollect string = "sync:sapo:collect"
)

type SyncSapoCollectPayload struct {
	Page       int
	Limit      int
	SapoDomain string
	StoreId    uint64
}

func NewSyncSapoCollectRecursiveTask(page, limit int, sapoDomain string, storeId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncSapoCollectPayload{Page: page, Limit: limit, SapoDomain: sapoDomain, StoreId: storeId})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncSapoCollect, payload), nil
}

func HandleSyncSapoCollectRecursiveTask(ctx context.Context, t *asynq.Task) error {
	var p SyncSapoCollectPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	count, err := sapo.SyncCollect(p.Page, p.Limit, p.SapoDomain, p.StoreId)
	if err != nil {
		return err
	}

	if count > 0 {
		payload, err := json.Marshal(SyncSapoCollectPayload{Page: p.Page + 1, Limit: p.Limit, SapoDomain: p.SapoDomain, StoreId: p.StoreId})
		if err != nil {
			return err
		}
		return HandleSyncSapoCollectRecursiveTask(ctx, asynq.NewTask(TypeSyncSapoCustomCollections, payload))
	}

	return nil
}
