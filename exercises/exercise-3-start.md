<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 3 - A real life matchmaking and leaderboard application using Go and Redis

## Goals
* Learn about [Redis Streams](https://redis.io/docs/data-types/streams-tutorial/)
* See some example usage of a Redis client library in application code
* Explore a matchmaking program
* Visualize Redis data in Grafana

## Structure of this exercise
* Background
* Setup
* Load up Grafana
* Start the matchmaking program

## Redis Streams

In this code, we'll be making use of a Redis data structure called Streams. Redis Streams act like append-only logs - items within a stream are immutable and can't be amended once they've been added to a stream. Redis generates a unique ID for each entry within a stream. Inside a stream, entries are similar to the structure of hashes, with a schemaless set of keys and values. Applications and services can be producers (add to a stream), or consumers (read from a stream). In the code, we'll be seeing examples of both producers and consumers.

For a deep dive into Redis Streams, check out [this tutorial](https://redis.io/docs/data-types/streams-tutorial/)

## Redis and Go

One of the great things about Redis is the number of client libraries available to developers. In this exercise, we'll be using Go, as well as [Redis client for Go](https://github.com/go-redis/redis), [RediSearch client for Go](https://github.com/RediSearch/redisearch-go) and the Redis [TimeSeries client for Go](https://github.com/go-redis/redis). 

If you're developing software but don't use Go, there are over 160 client libraries across over 50 different languages available. A selection is available [on the Redis website](https://redis.io/docs/clients/)

## How to do this in code

In exercise one, we covered some of the basic queries needed for matchmaking. Now let's see how we can integrate that into a (simple) program. The provisioned environments come with a preloaded dataset. We'll use this dataset to simulate incoming ticket requests from gamers, and this will take the form of a Redis Stream.
Then we'll match these players using Redis and RediSearch.

For this purpose we will need a number of things on Redis side: 

- An input stream, which contains tickets from players (or clients).
- A matchmaking service that listens to the stream and looks for potential matches.
- An output stream, for the matched games, to be given back to the client or downstream for further processing (e.g. to launch a game server). 
- We'll also emit metrics on another output stream so that we can visualize them within a dashboard.

## Diagram

![Diagram of matchmaker](/img/diagram.png)

## Setup

**When using a Redis Cloud database, update the following:**

Locate the [.env](./grafana/.env) file for docker compose and edit the file, substituting the values to match your Redis Cloud configuration:

```
#REDIS_SERVER=redis
#REDIS_PORT=6379
#REDIS_PASSWORD=

REDIS_SERVER=YOUR-REDIS-CLOUD-ENDPOINT
REDIS_PORT=YOUR-REDIS-CLOUD-PORT
REDIS_PASSWORD=YOUR-REDIS-CLOUD-PASSWORD
```

Edit Grafana's [datasources](./grafana/datasources/automatic.yml) to use Redis Cloud:

```
datasources:
  # - name: Redis
  #   type: redis-datasource
  #   access: proxy
  #   url: redis://redis:6379/
  - name: SimpleJson
    type: grafana-simple-json-datasource
    access: proxy
    url: http://leaderboard:5000/
  - name: Redis
    type: redis-datasource
    access: proxy
    url: redis://YOUR-REDIS-CLOUD-ENDPOINT:YOUR-REDIS-CLOUD-PORT/
    secureJsonData:
      password: YOUR-REDIS-CLOUD-PASSWORD
```


Start the docker containers:
```
cd exercises/grafana
docker-compose up
```

## Grafana 

Access the Grafana instance that will now be available at http://localhost:3000/. You'll be asked to set a username and password - use the default admin\admin credentials.

Within Grafana, there are two dashboards (Matchmaking and Leaderboards). For this exercise we're going to use the matchmaking dashboard to visualise the Redis data structures which are powering the matchmaking process.

Load up the matchmaking dashboard in your browser.

## Matchmaking code

For the next step, we're interested in `matchmaker.go`. There is also a separate stats process (`stats.go`) - but that is initialized by Docker Compose, so it's running and ready to receive events from the matchmaking service. 

In a terminal, run the matchmaker:
```
cd ../go
go run matchmaker.go
```

**When using a Redis Cloud database, pass command line arguments as follows:**

```
cd ../go
go run matchmaker.go --host YOUR-REDIS-CLOUD-ENDPOINT --port YOUR-REDIS-CLOUD-PORT --password YOUR-REDIS-CLOUD-PASSWORD
```

You will start to see data being populated within the panels of the Grafana Matchmaking dashboard. Let's take a moment to consider what's going on here:

* Incoming ticket requests are stored within a Redis Stream - the application is reading these events using [XREADGROUP](https://redis.io/commands/xreadgroup/).
* For each ticket request, a user lookup is performed on fields that are stored within a hash using [HMGET](https://redis.io/commands/hmget/).
* Attributes for a specific user are used to prepare and execute RediSearch queries (as seen in the examples in exercise one).
* The results from the query determine if a temporary data structure should be created or amended (a Redis Hash) until there's enough matching players for a given set of player skills and location. You'll note the use of [HMSET](https://redis.io/commands/hmset/).
* Matched players are added to another Redis Stream using [XADD](https://redis.io/commands/xadd/).
* The incoming tickets stream is cleaned up using [XDEL](https://redis.io/commands/xdel/) and hashes are removed using [UNLINK](https://redis.io/commands/unlink/) - a non blocking delete command.

## Visualizing the outputs

While the matchmaking process is running, take a look at the dashboard panels. Notice how the data is processed and changes over time.

The [Redis Datasource Plugin](https://grafana.com/grafana/plugins/redis-datasource/) is a great way to visualise Redis data structures.
