<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 2 - A real life matchmaking application using Go and Redis

## Goals
* Learn about [Redis Streams](https://redis.io/docs/data-types/streams-tutorial/)
* See some example usage of a Redis client library in application code
* Explore a matchmaking program
* Visualise Redis data in Grafana Insight

## Structure of this exercise
* Background
* Setup
* Start the matchmaking program
* Visualise the data in Grafana and RedisInsight

## Redis Streams

In this code, we'll be making use of a Redis data structure called Streams. Redis Streams act like append-only logs - items within a stream are immutable and can't be amended once they've been added to a stream. Redis generates a unique ID for each entry within a stream. Within the stream, entries are similar to structure of hashes, with a schema-less set of keys and values. Applications and services can be producers (add to a stream), or consumers (read from a stream). In the code, we'll be seeing examples of both producers and consumers.

For a deep dive into Redis Streams, check out [this tutorial](https://redis.io/docs/data-types/streams-tutorial/)

## Redis and Go

One of the great things about Redis is the number of client libraries available to developers. In this exercise, we'll be using Go, as well as [Redis client for Go](https://github.com/go-redis/redis), [RediSearch client for Go](https://github.com/RediSearch/redisearch-go) and the Redis [TimeSeries client for Go](https://github.com/go-redis/redis). If you're not using Go, don't worry, there are over 160 client libraries across over 50 different languages available. A selection is available [on the Redis website](https://redis.io/docs/clients/)

## How to do this in code

In exercise 1, we covered some of the basic queries needed for matchmaking. Now let's see how we can integrate that into a (simple) program. The provisioned environments come with a preloaded dataset. We'll use this dataset to simulate incoming ticket requests from gamers, and this will take the form of a Redis Stream.
Then we'll match these players using Redis and RediSearch.

For this purpose we will need a number of things on Redis side: 

- An input stream, which contains tickets from players (or clients).
- A match making service, that listens to the stream and looks for potential matches.
- An output stream, for the matched games, to be given back to the client or downstream for further processing (e.g. to launch a game server). 
- We'll also emit metrics on another output stream so that we can visualise them within a dashboard.

## Diagram

![Diagram of matchmaker](/img/diagram.png)

### Getting started with the code
First, let's fire up the Docker containers:
```
cd exercises/grafana
docker compose up
```
Second (in a separate tab/window), run the matchmaker:
```
cd ../go
go run matchmaker.go
```

### Exploring the code

By using a flow chart, let's take a look at what's going on in `matchmaker.go`:

![Diagram of matchmaker](/img/matchmaker-flow-diagram.png)

### Commands

[XREADGROUP](https://redis.io/commands/xreadgroup/)

[XADD](https://redis.io/commands/xadd/)

[XDEL](https://redis.io/commands/xdel/)

[HGET](https://redis.io/commands/hget/)

[HSET](https://redis.io/commands/hset/)

[SADD](https://redis.io/commands/sadd/)

[UNLINK](https://redis.io/commands/unlink/)

### Visualizing data
<tbd, explanation of grafana and dashboards>