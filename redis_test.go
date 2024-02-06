/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

// var (
// 	redisUrl      = os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT")
// 	redisPassword = os.Getenv("REDIS_PASSWORD")
// 	redisClient   = goredis.NewClient(&goredis.Options{Addr: redisUrl, Password: redisPassword, DB: 0})
// )

// func TestGetValuesFromRedisSet(t *testing.T) {

// 	os.Setenv("REDIS_PORT", "28015")
// 	os.Setenv("REDIS_PASSWORD", "Atlan7is")
// 	os.Setenv("REDIS_SERVER", "127.0.0.1")

// 	GetValueFromRedisByKey(redisClient, "pr-st-0-simulate-stagetime-1926523c5a")

// }

// var ctx = context.Background()

// // Name - student name
// type Name struct {
// 	First  string `json:"first,omitempty"`
// 	Middle string `json:"middle,omitempty"`
// 	Last   string `json:"last,omitempty"`
// }

// // Student - student object
// type Student struct {
// 	Name Name `json:"name,omitempty"`
// 	Rank int  `json:"rank,omitempty"`
// }

// func TestSetObjectToRedisJSON(t *testing.T) {

// 	// INITALIZE REDIS
// 	var addr = flag.String("Server", "localhost:6379", "Redis server address")

// 	redisJSONHandler := rejson.NewReJSONHandler()
// 	flag.Parse()

// 	redisClient := goredis.NewClient(&goredis.Options{Addr: *addr, DB: 0})

// 	redisJSONHandler.SetGoRedisClient(redisClient)

// 	student := Student{
// 		Name: Name{
// 			"Patrick",
// 			"Johannes",
// 			"Hermann",
// 		},
// 		Rank: 1,
// 	}

// 	SetObjectToRedisJSON(redisJSONHandler, student, "student456")

// }
