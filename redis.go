/*
Copyright © 2022 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

/*
name: AddValueToRedisSet
description: Function to add Value to a Key. If the key is non-existent it creates it as well.
return: The return value is a boolean that determines if the value already exists or is unique. (True = unique)
exampleUsage: |

	uniqueVal := AddValueToRedisSet(RedisClient, context, "NameGenerator", "DarthVader") //key = "NameGenerator"  value = "DarthVader"
*/
func AddValueToRedisSet(redisClient *redis.Client, setKey, value string) (isSetValueunique bool) {

	if redisClient.SAdd(setKey, value).Val() == 1 {
		isSetValueunique = true
	}

	return
}

/*
name: GetRandomValueFromRedis
description: Function that returns a random value from a specified key.
exampleUsage: |

	randomName := GetRandomValueFromRedis(RedisClient, context, "NameGenerator") // key = NameGenerator
*/
func GetRandomValueFromRedis(redisClient *redis.Client, ctx context.Context, key string) string {
	ranVal := redisClient.SRandMember(key).Val()
	return ranVal
}

/*
name: GetAllValuesFromRedis
description: Function that returns all values from an existing key.
exampleUsage: |

	rangeValues := GetAllValuesFromRedis(RedisClient, context, "NameGenerator") // key = NameGenerator
*/
func GetAllValuesFromRedis(redisClient *redis.Client, ctx context.Context, key string) []string {
	rvalues, err := redisClient.SMembers(key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return rvalues
}

/*
name: RemoveValueFromRedis
description: Function that removes a single specific value within a key
return: The return value is a boolean that determines if the value existed and got removed or did not exist. (True = existed)
exampleUsage: |

	removedValue := RemoveValueFromRedis(RedisClient, context, "NameGenerator", "DarthVader") // key = NameGenerator value="DarthVader"
*/
func RemoveValueFromRedis(redisClient *redis.Client, ctx context.Context, key string, member string) bool {
	remVal, err := redisClient.SRem(key, member).Result()
	if err != nil {
		fmt.Println(err)
	}
	if remVal == 1 {
		return true
	}
	return false
}
