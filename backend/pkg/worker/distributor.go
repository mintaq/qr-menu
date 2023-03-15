package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
)

var AsynqClient *asynq.Client

func CreateRedisClient() {
	redisConnURL, err := utils.ConnectionURLBuilder("redis")
	if err != nil {
		log.Panic(err.Error())
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: redisConnURL,
	}

	AsynqClient = asynq.NewClient(redisOpt)
	defer AsynqClient.Close()
}
