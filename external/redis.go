package external

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"

	"go-poc/utils/activity"
	"go-poc/utils/log"
)

func NewRedis() (*redis.Client, error) {
	ctx := activity.NewContext("init_redis")

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "can't ping redis db"))
		return nil, stacktrace.Propagate(err, "can't ping redis db")
	}

	return client, nil
}
