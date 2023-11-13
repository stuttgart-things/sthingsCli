/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"fmt"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	"github.com/stuttgart-things/redisqueue"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisJSON(redisJSONHandler *rejson.Handler, jsonKey string) (jsonObject []byte) {

	jsonObject, err := redigo.Bytes(redisJSONHandler.JSONGet(jsonKey, "."))

	fmt.Println(err)

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to JSONGet")
		return
	}

	return

}

func SetRedisJSON(redisJSONHandler *rejson.Handler, jsonObject interface{}, jsonKey string) {

	res, err := redisJSONHandler.JSONSet(jsonKey, ".", jsonObject)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}

	if res.(string) == "OK" {
		fmt.Printf("Success: %s\n", res)
	} else {
		fmt.Println("Failed to Set: ")
	}

}

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

func CheckRedisKV(connectionString, redisPassword, key, expectedValue string) (keyValueExists bool) {

	rdb := CreateRedisClient(connectionString, redisPassword)

	// CHECK IF KEY EXISTS IN REDIS
	fmt.Println("CHECKING IF KEY " + key + " EXISTS..")
	keyExists, err := rdb.Exists(context.TODO(), key).Result()
	if err != nil {
		panic(err)
	}

	// CHECK FOR VALUE/STATUS IN REDIS
	if keyExists == 1 {

		fmt.Println("KEY " + key + " EXISTS..CHECKING FOR IT'S VALUE")

		value, err := rdb.Get(context.TODO(), key).Result()
		if err != nil {
			panic(err)
		}

		if value == expectedValue {
			fmt.Println("STATUS", value)
			keyValueExists = true
		}

		fmt.Println("STATUS", value)

	} else {

		fmt.Println("KEY " + key + " DOES NOT EXIST)")
	}

	return
}

func EnqueueDataInRedisStreams(connectionString, redisPassword, stream string, values map[string]interface{}) (enqueue bool) {

	producer, err := redisqueue.NewProducerWithOptions(&redisqueue.ProducerOptions{
		MaxLen:               10000,
		ApproximateMaxLength: true,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     connectionString,
			Password: redisPassword,
			DB:       0,
		}),
	})

	if err != nil {
		panic(err)
	}

	redisStreamErr := producer.Enqueue(&redisqueue.Message{
		Stream: stream,
		Values: values,
	})

	if redisStreamErr != nil {
		panic(redisStreamErr)

	} else {
		enqueue = true
	}

	return
}
