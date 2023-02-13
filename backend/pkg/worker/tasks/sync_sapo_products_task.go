package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils/sapo"
)

const TypeSyncSapoProducts = "sync:sapo:products"

type SyncSapoProductsPayload struct {
	Page       int
	Limit      int
	SapoDomain string
	StoreId    uint64
}

func NewSyncSapoProductsRecursiveTask(page, limit int, sapoDomain string, storeId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncSapoProductsPayload{Page: page, Limit: limit, SapoDomain: sapoDomain, StoreId: storeId})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncSapoProducts, payload), nil
}

func HandleSyncSapoProductsRecursiveTask(ctx context.Context, t *asynq.Task) error {
	var p SyncSapoProductsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	count, err := sapo.SyncProducts(p.Page, p.Limit, p.SapoDomain, p.StoreId)
	if err != nil {
		return err
	}

	if count > 0 {
		payload, err := json.Marshal(SyncSapoProductsPayload{Page: p.Page + 1, Limit: p.Limit, SapoDomain: p.SapoDomain, StoreId: p.StoreId})
		if err != nil {
			return err
		}
		return HandleSyncSapoProductsRecursiveTask(ctx, asynq.NewTask(TypeSyncSapoProducts, payload))
	}

	return nil
}
