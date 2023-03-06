package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/kiotviet"
)

const TypeSyncKiotvietProducts = "sync:kiotviet:products"

type SyncKiotVietProductsPayload struct {
	UserId      uint64
	StoreId     uint64
	PageSize    int
	CurrentItem int
}

func NewSyncKiotVietProductsRecursiveTask(userId, storeId uint64, pageSize, currentItem int) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncKiotVietProductsPayload{UserId: userId, StoreId: storeId, PageSize: pageSize, CurrentItem: currentItem})
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

	lastCurrentItem, err := kiotviet.SyncProducts(p.UserId, p.StoreId, p.PageSize, p.CurrentItem)
	if err != nil {
		return err
	}

	if lastCurrentItem > 0 {
		payload, err := json.Marshal(SyncKiotVietProductsPayload{UserId: p.UserId, StoreId: p.StoreId, PageSize: p.PageSize, CurrentItem: lastCurrentItem})
		if err != nil {
			return err
		}
		return HandleSyncKiotvietProductsRecursiveTask(ctx, asynq.NewTask(TypeSyncKiotvietProducts, payload))
	}

	return nil
}
