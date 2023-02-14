package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils/kiotviet"
)

const TypeSyncKiotvietProducts = "sync:kiotviet:products"

type SyncKiotVietProductsPayload struct {
	Page  int
	Limit int
	Store string
}

func NewSyncKiotVietProductsRecursiveTask(page, limit int, store string) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncKiotVietProductsPayload{Page: page, Limit: limit, Store: store})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncKiotvietProducts, payload), nil
}

func HandleSyncKiotvietProductsRecursiveTask(ctx context.Context, t *asynq.Task) error {
	var p SyncKiotVietProductsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	count, err := kiotviet.SyncProducts(p.Page, p.Limit, p.Store)
	if err != nil {
		return err
	}

	if count > 0 {
		payload, err := json.Marshal(SyncKiotVietProductsPayload{Page: p.Page + 1, Limit: p.Limit, Store: p.Store})
		if err != nil {
			return err
		}
		return HandleSyncKiotvietProductsRecursiveTask(ctx, asynq.NewTask(TypeSyncKiotvietProducts, payload))
	}

	return nil
}
