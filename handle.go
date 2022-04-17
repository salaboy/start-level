package function

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"time"
)

var redisHost = os.Getenv("REDIS_HOST") // This should include the port which is most of the time 6379
var redisPassword = os.Getenv("REDIS_PASSWORD")

type GameTime struct{
	GameTimeId string
	SessionId string
	Level string
	Type string
	Time      time.Time
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	var gt GameTime
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&gt)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	gt.GameTimeId = "time-" + gt.SessionId
	gt.Type = "start"
	gt.Time = time.Now()

	gameTimeJson, err := json.Marshal(gt)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.RPush(gt.GameTimeId, string(gameTimeJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(res, string(gameTimeJson))
}
