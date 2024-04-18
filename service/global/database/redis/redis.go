package redis

import (
	"fmt"
	"simple-video-net/global/config"

	"github.com/go-redis/redis"
)

func ReturnsInstance() *redis.Client {
	var err error
	var redisConfig = config.Config.RConfig

	// Create a link
	Db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.IP, redisConfig.Port),
		Password: redisConfig.Password, // no password set
		DB:       0,                    // use default DB
	})
	_, err = Db.Ping().Result()
	if err != nil {
		fmt.Printf("redis connection failed.%v \n", err)
	}
	return Db

}
