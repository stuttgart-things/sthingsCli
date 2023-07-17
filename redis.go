/*
Copyright © 2022 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func AddValueToRedisSet(redisClient *redis.Client, setKey, value string) (isSetValueunique bool) {

	if redisClient.SAdd(context.TODO(), setKey, value).Val() == 1 {
		isSetValueunique = true
	}

	return
}

func GetValuesFromRedisSet(redisClient *redis.Client, setKey string) (values []string) {

	values = redisClient.SMembers(context.TODO(), setKey).Val()

	return
}

func GetRandomValueFromRedis(redisClient *redis.Client, ctx context.Context, key string) string {
	ranVal := redisClient.SRandMember(context.TODO(), key).Val()
	return ranVal
}

func GetAllValuesFromRedis(redisClient *redis.Client, ctx context.Context, key string) []string {
	rvalues, err := redisClient.SMembers(context.TODO(), key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return rvalues
}

func RemoveValueFromRedis(redisClient *redis.Client, ctx context.Context, key string, member string) bool {
	remVal, err := redisClient.SRem(context.TODO(), key, member).Result()
	if err != nil {
		fmt.Println(err)
	}
	if remVal == 1 {
		return true
	}
	return false
}

func CreateRedisClient(connectionString, redisPassword string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: redisPassword,
		DB:       0,
	})

	return
}
