package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	"github.com/pborman/getopt/v2"
)

func main() {

	var ctx = context.Background()

	helpFlag := getopt.BoolLong("help", 'h', "display help")
	//verbose := getopt.BoolLong("verbose", 'v', "verbose output")

	redisHost := getopt.StringLong("host", 's', "127.0.0.1", "Redis Host")
	redisPort := getopt.IntLong("port", 'p', 6379, "Redis Port")
	redisPassword := getopt.StringLong("password", 'a', "", "Redis Password")
	instream := getopt.StringLong("file", 'f', "tickets", "the name of the stream to read from")
	metricStream := getopt.StringLong("metrics-stream", 'm', "metrics", "where to stream metrics")
	backfillStream := getopt.StringLong("backfill-stream", 'b', "backfill", "Prefix to stream backfill notifications")
	leaderboardSet := getopt.StringLong("leaderboard-set", 'l', "LeaderBoards", "Which set we store leaderboards in")
	index := getopt.StringLong("index", 'i', "GameTix", "the city index to search from")
	getopt.Parse()

	if *helpFlag {
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	searchClient := redisearch.NewClient(
		fmt.Sprintf("%s:%d", *redisHost, *redisPort),
		*index,
	)
	// Create Search Index if not available

	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("ticket")).
		AddField(redisearch.NewTextField("pop")).
		AddField(redisearch.NewNumericField("mmr"))

	indexDefinition := redisearch.NewIndexDefinition().AddPrefix("ticket:")
	searchClient.CreateIndexWithIndexDefinition(schema, indexDefinition)

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

	redisClient.SAdd(ctx, *leaderboardSet, "LeaderBoard-Pop", "LeaderBoard-Users", "LeaderBoard-Groups")

	// Try to create a read group and it will fail if already present
	redisClient.XGroupCreateMkStream(
		ctx,
		*instream,
		fmt.Sprintf("Group-%s", *instream), "0").Err()

	for {
		res, _ := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    fmt.Sprintf("Group-%s", *instream),
			Consumer: fmt.Sprintf("Consumer-%s", *instream),
			Streams:  []string{*instream, ">"},
			Count:    10,
			Block:    time.Duration(50) * time.Millisecond,
		}).Result()
		for _, x := range res {
			for _, msg := range x.Messages {
				vals, err := redisClient.HMGet(ctx,
					fmt.Sprintf("user:%s", msg.Values["user"]),
					"pop", "mmr", "blacklist_tags", "group_tags", "secondary_group_tags", "play_style_tags", "experience").Result()
				if err != nil {
					log.Println(err)
				}
				mmr, _ := strconv.ParseFloat(vals[1].(string), 16)
				btags := msg.Values["user"].(string)
				if vals[2] != "" {
					btag := strings.Split(vals[2].(string), ",")
					btag = append(btag, msg.Values["user"].(string))
					btags = strings.Join(btag, "|")
				}
				xp, _ := strconv.Atoi(vals[6].(string))
				xp = xp * (rand.Intn(1000) + 300)
				redisClient.XAdd(ctx, &redis.XAddArgs{Stream: *metricStream, Values: map[string]interface{}{
					"stat":             msg.Values["user"].(string),
					"type":             "gauge",
					"score":            xp,
					"leaderboard":      2,
					"leaderboard_name": "LeaderBoard-Users",
				}})
				searchQuery := fmt.Sprintf("@pop:%s @mmr:[%d %d] -@user:(%s)", vals[0], int(mmr*0.975), int(mmr*1.025), btags)

				// add optional group search params
				if vals[3] != "" {
					searchQuery += fmt.Sprintf(" ~@group_tags:{%s}", strings.ReplaceAll(vals[3].(string), ",", "|"))
					/* for _, x := range strings.Split(vals[3].(string), ",") {
						redisClient.XAdd(ctx, &redis.XAddArgs{Stream: *metricStream, Values: map[string]interface{}{
							"stat":             x,
							"type":             "counter",
							"leaderboard":      2,
							"leaderboard_name": "LeaderBoard-Groups",
						}})
					} */
				}
				if vals[4] != "" {
					searchQuery += fmt.Sprintf(" ~@secondary_group_tags:{%s}", strings.ReplaceAll(vals[4].(string), ",", "|"))
				}
				if vals[5] != "" {
					searchQuery += fmt.Sprintf(" ~@play_style_tags:{%s}", strings.ReplaceAll(vals[5].(string), ",", "|"))
				}

				docs, _, serr := searchClient.Search(redisearch.NewQuery(searchQuery).Limit(0, 3))
				if serr != nil {
					log.Println("SEARCHERROR: ", serr, searchQuery)
				}
				if len(docs) < 1 {
					redisClient.HMSet(
						ctx,
						fmt.Sprintf("ticket:%s:%s", msg.ID, msg.Values["user"]),
						"ticket", msg.ID,
						"user", msg.Values["user"],
						"pop", vals[0],
						"mmr", vals[1],
						"group_tags", vals[3],
						"secondary_group_tags", vals[4],
						"play_style_tags", vals[5],
					)
				} else {
					tix := docs[0].Properties["ticket"]
					redisClient.HMSet(
						ctx,
						fmt.Sprintf("ticket:%s:%s", tix, msg.Values["user"]),
						"ticket", tix,
						"user", msg.Values["user"],
						"pop", vals[0],
						"mmr", vals[1],
						"group_tags", vals[3],
						"secondary_group_tags", vals[4],
						"play_style_tags", vals[5],
					)

					if len(docs) == 3 {
						group := []string{msg.Values["user"].(string)}
						for _, doc := range docs {
							group = append(group, doc.Properties["user"].(string))
						}
						redisClient.Unlink(ctx,
							fmt.Sprintf("ticket:%s:%s", tix, msg.Values["user"]),
							docs[0].Id,
							docs[1].Id,
							docs[2].Id,
						)

						//if *verbose {
						log.Printf("Grouped!! %+v\n", group)
						//}
						// Send match to backfill service
						redisClient.XAdd(ctx, &redis.XAddArgs{
							Stream: fmt.Sprintf("%s-%s", *backfillStream, vals[0]),
							Values: map[string]interface{}{
								"pop":   vals[0],
								"users": strings.Join(group, ","),
							},
						})
					}
				}
				redisClient.XDel(ctx, *instream, msg.ID)
				/* redisClient.XAdd(ctx, &redis.XAddArgs{Stream: *metricStream, Values: map[string]interface{}{
					"stat":        "match",
					"type":        "counter",
					"leaderboard": 0,
				}})
				redisClient.XAdd(ctx, &redis.XAddArgs{Stream: *metricStream, Values: map[string]interface{}{
					"stat":             fmt.Sprintf("%s-PoP", vals[0]),
					"type":             "counter",
					"leaderboard":      1,
					"leaderboard_name": "LeaderBoard-Pop",
				}}) */
			}
		}
	}
}
