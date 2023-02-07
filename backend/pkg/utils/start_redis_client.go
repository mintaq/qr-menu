package utils

import (
	"log"

	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

func CreateRedisClient() {
	redisConnURL, err := ConnectionURLBuilder("redis")
	if err != nil {
		log.Panic(err.Error())
	}

	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisConnURL})
}
