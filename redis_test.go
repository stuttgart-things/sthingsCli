/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

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
