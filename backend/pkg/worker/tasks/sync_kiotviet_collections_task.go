package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/kiotviet"
)

const TypeSyncKiotvietCollections = "sync:kiotviet:collections"

type SyncKiotVietCollectionsPayload struct {
	UserId      uint64
	StoreId     uint64
	PageSize    int
	CurrentItem int
}

func NewSyncKiotVietCollectionsRecursiveTask(userId, storeId uint64, pageSize, currentItem int) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncKiotVietCollectionsPayload{UserId: userId, StoreId: storeId, PageSize: pageSize, CurrentItem: currentItem})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncKiotvietCollections, payload), nil
}

func HandleSyncKiotvietCollectionsRecursiveTask(ctx context.Context, t *asynq.Task) error {
	var p SyncKiotVietCollectionsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	lastCurrentItem, err := kiotviet.SyncCollections(p.UserId, p.StoreId, p.PageSize, p.CurrentItem)
	if err != nil {
		return err
	}

	if lastCurrentItem > 0 {
		payload, err := json.Marshal(SyncKiotVietCollectionsPayload{UserId: p.UserId, StoreId: p.StoreId, PageSize: p.PageSize, CurrentItem: lastCurrentItem})
		if err != nil {
			return err
		}
		return HandleSyncKiotvietCollectionsRecursiveTask(ctx, asynq.NewTask(TypeSyncKiotvietCollections, payload))
	}

	return nil
}
