package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
)

// GetCartData retrieves the cart data stored in Redis based on a given key.
// If the key does not exist in Redis, a new cart is created with the provided cart token.
// If the cart token is already taken, an error is returned.
func GetCartData(key, cartToken string) (*models.Cart, error) {
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	cart := &models.Cart{}
	if err == redis.Nil {
		cart.CartToken = cartToken
		err = SetCartData(key, cart, utils.GetRedisCartDuration())
		if err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte(val), cart); err != nil {
			return nil, err
		}

		if cart.CartToken != "" && cart.CartToken != cartToken {
			return nil, errors.New("table is already taken")
		}
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
