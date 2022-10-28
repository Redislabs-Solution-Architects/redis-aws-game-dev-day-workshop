<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 2 - A real life matchmaking application using Go and Redis

## Goals
* See some example usage of Redis client libraries in application code
* Explore a matchmaking program
* Learn about [Redis Streams](https://redis.io/docs/data-types/streams-tutorial/)
* Visualise Redis data in Grafana Insight

## Structure of this exercise
* Setup
* Clone the source code repository
* Start the matchmaking program
* Visualise the data in Grafana and RedisInsight

## Redis and Go

One of the great things about Redis is the number of client libraries available to developers. In this exercise, we'll be using Go, as well as [Redis client for Go](https://github.com/go-redis/redis), [RediSearch client for Go](https://github.com/RediSearch/redisearch-go) and the Redis [TimeSeries client for Go](https://github.com/go-redis/redis).

### How to do this all in code

Now that we've covered some of the basic queries needed for matchmaking, let's see how we can integrate that into a (simple) program. The provision environments come with a preloaded dataset of users. We'll use this dataset to generate a set of tickets/match making requests, and then we'll match them using Redis and RediSearch.

For this purpose we will need a number of things on Redis side: 

- an input Stream, which players/clients would use to put their requests/tickets on.
- A MatchMaker, that listens to the stream and looks for potential matches. 
- And an output Stream, for the matched games, to be given back to the client or downstream for further processing (e.g. looking for an available game server if it's not P2P, etc., etc.) 
- We'll also emit metrics on another output Stream.

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

Let's take a look at what's going on in `matchmaker.go`:

- We're working with three streams - one for tickets, one for backfills, and another for metrics.
- You'll notice a Redis set: 'LeaderBoards'
- Firstly, two Redis clients are instantiated (a Redis client, and RediSearch client).
- In order to perform RedisSearch queries, a schema needs to be set up - in this case, for 'ticket', 'pop' and 'mmr'.
- An input stream ('tickets'), is read using XREADGROUP - this means you can read from a stream with a consumer group.
```
XReadGroup(ctx context.Context, a *XReadGroupArgs)
```
- Next, for each message in the stream, a lookup is performed on the Redis database. You'll recognise this datatype from exercise 1 - we are retrieving a hash (a user) from Redis.
```
HMGet(ctx context.Context, key string, fields ...string)
```
- Once we have a user, variables are set for matchmaking rank, blocklist tags and experience score.
- At this point, the user is added to the User Leaderboard ('LeaderBoard-Users') by adding the user to the metric stream:
```
XAdd(ctx context.Context, a *XAddArgs)
```
- Now you'll notice the definition of a RediSearch query:
```
searchQuery := fmt.Sprintf("@pop:%s @mmr:[%d %d] -@user:(%s)", vals[0], int(mmr*0.975), int(mmr*1.025), btags)
```
- Next, some filtering on the message from the ticket stream is performed to see if the user belongs to a group.  If they are, the Groups Leaderboard is incremented by adding another message to the metric stream.
- Some more filtering is performed to amend the query.
- The RedisSearch query is now performed. If no results are returned, a ticket (in the form of a Redis hash) is created (using `redisClient.HMSet()`). Otherwise, it's updated for the matching user.
```
searchClient.Search(redisearch.NewQuery(searchQuery)
```
- The next step is where the actual matmaching occurs. If the RediSearch query returns three users ready to play, a fourth is added to the group, and this data is put on the backfill stream. Tickets that have been persisted as hashes get unlinked (a non-blocking Redis delete command) at this point.
- Some cleanup is performed by deleting the message from the ticket stream.
```
XDel(ctx context.Context, stream string, ids ...string)
```
- Finally, the matchmaking and PoP statistics are updated.

### Commands

[XREADGROUP](https://redis.io/commands/xreadgroup/)

[XADD](https://redis.io/commands/xadd/)

[XDEL](https://redis.io/commands/xdel/)

[HGET](https://redis.io/commands/hget/)

[HSET](https://redis.io/commands/hset/)

[SADD](https://redis.io/commands/sadd/)

[UNLINK](https://redis.io/commands/unlink/)