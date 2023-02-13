package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils/sapo"
)

const (
	TypeSyncSapoCustomCollections string = "sync:sapo:custom_collections"
	TypeSyncSapoSmartCollections  string = "sync:sapo:smart_collections"
)

type SyncSapoCollectionsPayload struct {
	Page  int
	Limit int
	Store string
}

func NewSyncSapoCustomCollectionsRecursiveTask(page, limit int, store string) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncSapoCollectionsPayload{Page: page, Limit: limit, Store: store})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncSapoCustomCollections, payload), nil
}

func HandleSyncSapoCustomCollectionsRecursiveTask(ctx context.Context, t *asynq.Task) error {
	var p SyncSapoCollectionsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	count, err := sapo.SyncCustomCollections(p.Page, p.Limit, p.Store)
	if err != nil {
		return err
	}

	if count > 0 {
		payload, err := json.Marshal(SyncSapoCollectionsPayload{Page: p.Page + 1, Limit: p.Limit, Store: p.Store})
		if err != nil {
			return err
		}
		return HandleSyncSapoCustomCollectionsRecursiveTask(ctx, asynq.NewTask(TypeSyncSapoCustomCollections, payload))
	}

	return nil
}

func NewSyncSapoSmartCollectionsRecursiveTask(page, limit int, store string) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncSapoCollectionsPayload{Page: page, Limit: limit, Store: store})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncSapoSmartCollections, payload), nil
}

func HandleSyncSapoSmartCollectionsRecursiveTask(ctx context.Context, t *asynq.Task) error {
	var p SyncSapoCollectionsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	count, err := sapo.SyncSmartCollections(p.Page, p.Limit, p.Store)
	if err != nil {
		return err
	}

	if count > 0 {
		payload, err := json.Marshal(SyncSapoCollectionsPayload{Page: p.Page + 1, Limit: p.Limit, Store: p.Store})
		if err != nil {
			return err
		}
		return HandleSyncSapoSmartCollectionsRecursiveTask(ctx, asynq.NewTask(TypeSyncSapoSmartCollections, payload))
	}

	return nil
}
