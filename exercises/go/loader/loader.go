package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pborman/getopt/v2"
)

func main() {

	var ctx = context.Background()

	helpFlag := getopt.BoolLong("help", 'h', "display help")

	redisHost := getopt.StringLong("host", 's', "127.0.0.1", "Redis Host")
	redisPort := getopt.IntLong("port", 'p', 6379, "Redis Port")
	redisPassword := getopt.StringLong("password", 'a', "", "Redis Password")
	pipesize := getopt.IntLong("pipeline", 'q', 10, "Pipeline size")
	stream := getopt.StringLong("stream", 'x', "tickets", "the name of the streams created")
	infile := getopt.StringLong("file", 'f', "", "the name of the csv to read from")
	getopt.Parse()

	if *helpFlag || *infile == "" {
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	csvFile, err := os.Open(*infile)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", *redisHost, *redisPort),
		Password:        *redisPassword,
		DB:              0,
		MinIdleConns:    2,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 5000 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      10, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	pipe := client.Pipeline()
	p := 1

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if p > 1 {

			pipe.XAdd(ctx, &redis.XAddArgs{Stream: *stream, Values: map[string]interface{}{
				"user": row[1],
			}})
			if p%*pipesize == 0 {
				_, pipeerr := pipe.Exec(ctx)
				if pipeerr != nil {
					log.Printf("ERROR flushing pipeline at %d end: %s\n", p, pipeerr)
				}
			}
		}
		p += 1

	}
	_, finalerr := pipe.Exec(ctx)
	if finalerr != nil {
		log.Printf("ERROR flushing pipeline at end: %s\n", finalerr)
	}

}
