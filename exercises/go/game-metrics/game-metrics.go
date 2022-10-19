package main

import (
	"context"
	"math"
	"math/rand"
	"time"

	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	redis "github.com/go-redis/redis/v8"
)

func main() {

	var ctx = context.Background()

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var timeSeriesClient = redistimeseries.NewClient("localhost:6379", "nohelp", nil)

	now := (time.Now().Unix() / 86400) * 86400

	for i := 364; i >= 0; i-- {

		timestamp := 1000 * (now - int64((86400 * i)))

		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 10
		random := rand.Intn(max-min) + min

		// Set some ARPU numbers
		arpu := 1 + float64(random) + math.Min(1.0, 1.2) + (float64(365-i) * 0.25)
		timeSeriesClient.Add("ARPU", timestamp, arpu)
		redisClient.Set(ctx, "ARPU_TODAY", arpu, 0)

		// daily active users
		dau := 10000000 + (365-i)*8220
		timeSeriesClient.Add("DAU", timestamp, float64(dau))
		redisClient.Set(ctx, "DAU_TODAY", dau, 0)

		// build some load times
		var login_time float64 = 0
		var daily_purchases float64 = 0

		if i > 90 {
			login_time = 10 + 3*rand.Float64()
			daily_purchases = float64(dau) * (arpu / 30)
		} else {
			login_time = 2 + rand.Float64()
			daily_purchases = float64(dau) * (arpu / 30)
		}

		app_time := 21 + 2*rand.Float64()

		timeSeriesClient.Add("LOAD_LOGIN", timestamp, (float64(login_time)))
		redisClient.Set(ctx, "LOGIN_TODAY", login_time, 0)

		timeSeriesClient.Add("LOAD_APP", timestamp, (float64(app_time)))
		redisClient.Set(ctx, "APP_TODAY", app_time, 0)

		timeSeriesClient.Add("DAILY_PURCHASES", timestamp, (float64(daily_purchases)))
		redisClient.Set(ctx, "REVENUE_TODAY", daily_purchases, 0)
	}

}
