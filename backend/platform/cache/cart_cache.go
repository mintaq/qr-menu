package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
)

func GetCartData(key string) (*models.Cart, error) {
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key %s does not exist", key)
		}
		return nil, err
	}

	cart := new(models.Cart)
	if err := json.Unmarshal([]byte(val), cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func SetCartData(key string, cart *models.Cart, expireTime time.Duration) error {
	dataStr, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	if err := RedisClient.Set(context.Background(), key, dataStr, expireTime).Err(); err != nil {
		return err
	}

	return nil
}
