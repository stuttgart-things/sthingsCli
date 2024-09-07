/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/RediSearch/redisearch-go/redisearch"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	"github.com/stuttgart-things/redisqueue"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func GetRedisJSON(redisJSONHandler *rejson.Handler, jsonKey string) (jsonObject []byte, jsonExists bool) {

	jsonObject, err := redigo.Bytes(redisJSONHandler.JSONGet(jsonKey, "."))

	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to JSONGet")
		return

	} else {
		jsonExists = true
	}

	return

}

func SetRedisJSON(redisJSONHandler *rejson.Handler, jsonObject interface{}, jsonKey string) {

	res, err := redisJSONHandler.JSONSet(jsonKey, ".", jsonObject)

	if err != nil {
		log.Fatalf("FAILED TO JSONSET")
		return
	}

	if res.(string) == "OK" {
		fmt.Printf("SUCCESS: %s\n", res)
	} else {
		fmt.Println("FAILED TO SET: ")
	}

}

func DeleteRedisSet(redisClient *redis.Client, setKey string) (isSetDeleted bool) {

	err := redisClient.Del(context.TODO(), setKey).Err()
	if err != nil {
		panic(err)
	} else {
		isSetDeleted = true
	}

	return
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

func GetValueFromRedisByKey(redisClient *redis.Client, key string) (value string) {

	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}

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

// CHECK IF REDISEARCH-INDEX EXISTS
func CheckIfRedisSearchIndexExists(client *redisearch.Client) (bool, error) {
	_, err := client.Info()

	if err != nil {
		if err.Error() == "Unknown Index name" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CREATE A REDIS CONNECTION POOL
func CreateRedisConnectionPool(redisHost, redisPassword string) (pool *redigo.Pool) {

	pool = &redigo.Pool{
		MaxIdle:   10,
		MaxActive: 10,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", redisPassword); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	return

}

// CREATE REDISEARCH INDEX
func CreateRedisSearchIndex(client *redisearch.Client, schema *redisearch.Schema) {

	if err := client.CreateIndex(schema); err != nil {
		log.Fatalf("Could not create index: %v", err)
	}

}

// DROP REDISEARCH INDEX
func DropRedisSearchIndex(client *redisearch.Client) {

	client.Drop()

}

func SearchQuery(client *redisearch.Client, query *redisearch.Query) {

	// SEARCH FOR DOCUMENTS
	docs, total, err := client.Search(query)
	if err != nil {
		log.Fatalf("Could not search for documents: %v", err)
	}

	fmt.Printf("Found %d documents\n", total)
	for _, doc := range docs {
		fmt.Printf("Document: %v\n", doc)
	}
}
