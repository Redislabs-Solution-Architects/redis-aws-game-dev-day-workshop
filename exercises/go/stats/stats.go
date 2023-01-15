package main

import (
	"context"
	"fmt"
	"strconv"

	"log"
	"os"
	"strings"
	"time"

	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/go-redis/redis/v8"
	"github.com/pborman/getopt/v2"
)

func main() {

	var ctx = context.Background()

	helpFlag := getopt.BoolLong("help", 'h', "display help")
	redisHost := getopt.StringLong("host", 's', "127.0.0.1", "Redis Host")
	redisPort := getopt.IntLong("port", 'p', 6379, "Redis Port")
	redisPassword := getopt.StringLong("password", 'a', "", "Redis Password")
	metricStream := getopt.StringLong("metrics-stream", 'm', "metrics", "where to stream metrics")
	getopt.Parse()

	if *helpFlag {
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	var tsPassword *string
	if *redisPassword != "" {
		tsPassword = redisPassword
	}

	// Timeseries Client
	timeClient := redistimeseries.NewClient(
		fmt.Sprintf("%s:%d", *redisHost, *redisPort),
		"nohelp",
		tsPassword,
	)

	// Stream client
	redisClient := redis.NewClient(&redis.Options{
		Password:        *redisPassword,
		Addr:            fmt.Sprintf("%s:%d", *redisHost, *redisPort),
		DB:              0,
		MinIdleConns:    1,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 5000 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      10, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	// Try to create a read group and it will fail if already present
	redisClient.XGroupCreateMkStream(
		ctx,
		*metricStream,
		fmt.Sprintf("Group-%s", *metricStream), "0").Err()

	for {
		res, _ := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    fmt.Sprintf("Group-%s", *metricStream),
			Consumer: fmt.Sprintf("Consumer-%s", *metricStream),
			Streams:  []string{*metricStream, ">"},
			Count:    10,
			Block:    time.Duration(50) * time.Millisecond,
		}).Result()
		for _, x := range res {
			for _, msg := range x.Messages {
				//fmt.Println(msg)
				ts, _ := strconv.ParseInt(strings.Split(msg.ID, "-")[0], 10, 64)
				switch statType := msg.Values["type"].(string); statType {
				case "counter":
					// if we set the stat name to 0 then we don't record a timeseries
					if msg.Values["leaderboard"].(string) != "2" {
						//update counters here
						timeClient.IncrBy(
							msg.Values["stat"].(string),
							ts,
							1.0,
							redistimeseries.CreateOptions{},
						)
					}
					if msg.Values["leaderboard"].(string) != "0" {
						if msg.Values["type"].(string) == "counter" {
							redisClient.ZIncr(
								ctx,
								msg.Values["leaderboard_name"].(string),
								&redis.Z{
									Member: msg.Values["stat"].(string),
									Score:  1,
								},
							)
						}
					}
				case "gauge":
					if msg.Values["leaderboard"].(string) != "0" {
						sc, _ := strconv.ParseFloat(msg.Values["score"].(string), 64)
						redisClient.ZIncr(
							ctx,
							msg.Values["leaderboard_name"].(string),
							&redis.Z{
								Member: msg.Values["stat"].(string),
								Score:  sc,
							},
						)
					}
				default:
					log.Printf("Unknown stat type %s\n", statType)
				}
				redisClient.XDel(ctx, *metricStream, msg.ID)

			}
		}
	}
}
