<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 1 - Matchmaking
Hello and thank you for joining this hands-on lab! We're happy to have you join, so welcome! This first exercise aims to get you up and running and introduce you to some of the core concepts and datatypes of Redis that we will be using to build a very fast Matchmaking engine. We'll also introduce you to the data type that powers many Leaderboards across the world. If you're stuck, don't dwell too long on it. Instead, check out the hints and the solution [over here](exercise-1-solution.md) or reach out to one of the instructors, who will be quite happy to help! Also, be sure to check out Volkan Civelek's excellent e-book on Matchmaking at <tbd>.

## Goals

* Learn some of the basic Redis commands by using the [Redis CLI](https://redis.io/topics/rediscli)
* Learn about some of the different [data types in Redis](https://redis.io/topics/data-types-intro)
* Learn about the [Hash](https://redis.io/docs/data-types/hashes/) data type and how we will use it to manage Tickets
* Learn how to build a Leaderboard using a [Sorted Set](https://redis.io/docs/data-types/sorted-sets/)
* Learn how to use [RediSearch](https://redis.io/docs/stack/search/) to execute fast queries on Redis Hashes or [JSON](https://redis.io/docs/stack/json/) documents to build a fast matchmaking engine.

## Structure of this exercise
1. Setup
1. Basic Redis interaction
1. Matchmaking basics and how you can use Redis and RediSearch to perform this fast and at high througput
1. How to do this in code

### Starting Redis and opening the CLI
* In a terminal Window, open the Redis CLI and connect to your provided Redis Enterprise database by typing:
```
redis-cli -h <hostname> -p <port>
```

By default (when executed without options), the Redis CLI will connect to 127.0.0.1 and port 6379. If you want to use a different hostname and port, you can use the `-h` and `-p` options to specify a diferent hostname and port.

### Authenticating
* The provided databases are secured using a password, so before we can do anything we need to provided our authentication. From the Redis CLI:
```
auth <password>
```
and the server should respond with 'OK'.

### Strings and hashes
* The databases that are provided come preloaded with the data required for this exercise, but let's practice some basics first. Let's see if we can store something in Redis and retrieve it again and start with storing/retrieving a basic String key/value pair:
* Add a String key/value pair to Redis with a key of `hello` and a value of `aws`:
```
set hello aws
```
* Retrieve the key/value pair from Redis
```
get hello
```
Congratulations, you are now a Redis veteran! Of course, `get` and `set` are not the only two commands available for dealing with Strings in Redis. You can find out about all of the other commands for Strings (and all the other datatypes) by visiting the [Redis documentation](https://redis.io/commands/#string)

Let's check out some of the other data structures in Redis as well, as we will be using some of these during the exercises. Let's start with a Hash.

Hashes are typically used to store flat structures with multiple attributes, such as records or structs. A very common use case is to store player profiles and game sessions in Redis so you can keep your own application stateless and add/remove instances/pods as needed while still being able to retrieve and modify state from a logged in player very quickly. In our exercises we will be using Hashes to store our users/players, which we can then query using RediSearch.

We will use matching criteria such as `MMR` (Matchmaking Ranking), `experience` (# of matches played), `location` (to match with players close geographically) and a few others.

Note that with Redis Enterprise, you can also setup Active-Active across regions (and even clouds if you so desire) so your players always have fast, local access to their sessions/profiles/tickets/etc. while the CRDT synchronisation makes sure that same data is also available in other regions.

The provided database comes with many users pre-loaded, so let's query a few and find out what Hashes are all about.

```
* Retrieve the Hash and all its attribute key/value pairs from Redis
```
hgetall user:ygreen6478
```
* Retrieve a single attribute from the Hash
```
hget user:ygreen6478 mmr
```

Did you notice that Redis has different commands for different data structures? For instance, a `get` command will not work on a hash, but it will work on all other data structures. And an `hget` will work on a hash but not on another data structure. If you try this you will get a 'WRONGTYPE' error. Don't know the type of a certain key? You can ask Redis the type of the `user:lars` key by typing:
```
type user:ygreen6478
```
and Redis will tell you that it's a hash. Alternatively, you can type:
```
type hello
```
To find out that this key is a String type.

If you want to learn more about hashes then check out all the available commands at the [Redis Documentation](https://redis.io/commands#hash). You can use the drop down list on that [same page](https://redis.io/commands) to select a different data structure to find out what commands apply to that particular data structure.

Now that you have seen a few of the most commonly used Redis data structures in action, feel free to take a look at what other data types and command are available in the [Redis documentation](https://redis.io/topics/data-types)

### Modules
 But wait, there's more! While Redis has a very powerful set of data structures straight out of the box, it's always possible that you find yourself in a situation where these are not enough to cover your use cases. In those situations Redis Modules comes to the rescue! Modules allow you to extend Redis data structures and features in a modular way; you only add those modules to a specific database that you want and no more. This keeps your Redis setup fast and lean. There are many community modules available and Redis provides several modules itself as well. If that's not enough you can also write your own. The environment that was provisioned for you has a few modules bundled with it already and we'll be using them extensively in this exercise. For more information on modules, see the [Redis Modules Hub](https://redis.com/community/redis-modules-hub/) and the [Redis Enterprise Modules](https://redis.com/redis-enterprise/modules/) pages.

### RediSearch
Let's start with RediSearch, a full-text search module for Redis. Retrieving keys by their primary value is fine for many use cases, but what if I have data in the cache that needs to be searched or aggregated? E.g. a product model, or stores, or transactions by a certain vendor, etc. Typically this is where a key/value model starts to show its limitations, e.g. having to transfer large chunks of data to the client before being able to search or otherwise process it. RediSearch to the rescue! RediSearch adds search capabilities to Redis as well as a lot more, and it does so right where your data lives! This means you can process much larger volumes of data and much faster as well.

See the [RediSearch documentation](https://docs.redis.com/latest/modules/redisearch/) for more detailed information on this module, or checkout the [GitHub repo](https://github.com/RediSearch/RediSearch).

In the context of matchmaking, RediSearch can be used to quickly fill games with matching players by doing the searching and matching inside Redis, where your data lives. This allows for much faster in-memory matching and will give you the highest possible throughput at the lowest possible latency with almost zero interaction on the client. Let's take a look at some of the query's that make this possible.

* First, we'll create a search index on Hash structures matching a certain prefix:
```
FT.CREATE GameTix ON HASH PREFIX 1 user: SCHEMA username TEXT mmr NUMERIC SORTABLE experience NUMERIC location GEO play_style_tags TAG blacklist_tags TAG group_tags TAG secondary_group_tags TAG pop TEXT SORTABLE
```

This command is a little bit more elaborate than the previous ones, so let's explore it in detail a bit more. We're creating an index called `Game-x` on the `hash` datastructure with the prefix `user:` (remember that we created a Hash earlier that had the key `user:lars`?) and we define the schema to be on the `username` field, which we define as a `TEXT` field. Note that you can also do really cool things such as phonetic search, but that's not needed for this exercise.

Furthermore we also add the `mmr`, `experience`, `location`, `play_style_tags`, `blacklist_tags`, `group_tags`, `secondary_group_tags` and `pop` fields to the index. Note that you can matchmake on one, more or all items, depending on your use case. But keep in mind that the more attributes you match on, the narrower your search is going to be and it will result in fewer results. So there's always the need to balance quality with quantity, but this may very will depend on the specific player.

Now let's perform some searches that are relevant in a matchmaking context:

Let's say we want to match players that are nearby geographically speaking, we need to match their locations with a certain distance, both of each other and where the game server may live (if any).

We'll start with the simplest of searches:
```
FT.SEARCH GameTix "*" 
```
This will return all items that are indexed. Of course this is not very useful to us, so maybe we should limit our search to items that are in a certain geography, e.g. players on the closest server:

```
FT.SEARCH GameTix "@pop:Auckland LIMIT 0 4" 
```

This limits our search already, but perhaps it's also a good idea to limit items to a certain MMR (Matchmaking Rating):
```
FT.SEARCH GameTix "@pop:Auckland @mmr:[2616 2817]"
```

Notice how we specify a range in MMR, allowing us to match players close, but not exactly at the same MMR as our own. Now let's limit it further by excluding players that we blocked and including players with a certain play style:

```
FT.SEARCH GameTix "@pop:Auckland @mmr:[2616 2817] -@user:(jabbott|hughesivan) ~@play_style_tags:{high_mobile}"
```

The `-` prefix before `@user` means we want to exclude any players that match on this criteria. The `~` before `play_style_tags` means we want to include players with this playstyle, but it's optional. Players that do match on playstyle will get scored higher in the results than players that do not match. This allows you to be more flexible in matchmaking and include more and more parameters without actually limiting the results, but rather prioritize them, as below:

```
FT.SEARCH GameTix "@pop:Auckland @mmr:[2616 2817] -@user:(jabbott|hughesivan) ~@play_style_tags:{high_mobile} ~@group_tags:{lightcyan_squad} ~@secondary_group_tags:{olivedrab_posse}" WITHSCORES
```

Besides being more complex, this query also shows the scores of the individual items. So how are scores calculated? We can executed the same query, but with one additional parameter:

```
FT.SEARCH GameTix "@pop:Auckland @mmr:[2616 2817] -@user:(jabbott|hughesivan) ~@play_style_tags:{high_mobile} ~@group_tags:{lightcyan_squad} ~@secondary_group_tags:{olivedrab_posse}" WITHSCORES EXPLAINSCORE
```

By default, the scoring function is [TFIDF](https://en.wikipedia.org/wiki/Tf%E2%80%93idf), but there are others available, see the [documentation](https://redis.io/docs/stack/search/reference/scoring/) for more info. You can also add your own scoring function!

We can also apply aggregate transformations on search result. E.g. what if we want to find the closest city/server to the player? E.g. we can go over all the cities/server in the data set and then apply a `geodistance()` transformation, e.g. sorting it according to distance to the players long/lat. Try running the query below:

* ### Find the closest network edge POP to the player (based on long/lat)

```
FT.AGGREGATE cities '*' LOAD 2 location city APPLY "geodistance(@location, -80.1401415,25.8102415)" as dist SORTBY 2 @dist ASC LIMIT 0 1
```

The output of this query could then be used

For a full overview of the RediSearch Query syntax, click [here](https://redis.io/docs/stack/search/reference/query_syntax/).

### Explanation and background
When matchmaking, we want to optimise the quality of the the matched game, which often boils down to a number of things:

1. Match with people of similar skill level
1. Match with people close to your geographical location
1. Match with people that match/complement your play style
1. Do NOT match with people that have been blocked by you
1. Shortest wait time as possible

To ensure the greatest chance of having a game of sufficient quality, the pool of people to match must be as large as possible. Which means that if your matchmaking algorithm is slow, you will be able to handle less matches/second and are likely facing increased latency (meaning more wait time for everyone). This is why you'd rather not do this on the application side, as it requires a transfer of data and processing power on the consuming side. Do this on the server, right where your data lives; which is why Redis + RediSearch is such a great idea for this scenario. It has no problem to handle a pool of several millions of users/tickets and provide quality matches. E.g. think about all the extra work needed and potentially extra wait time if you would need to split that pool up into smaller sections to be able to process it.

If you want to learn more about all the functionality of the RediSearch module, please check the [RediSearch documentation](https://oss.redis.com/redisearch/).

### How to do this all in code

Now that we've covered some of the basic queries needed for matchmaking, let's see how we can integrate that into a (simple) program. The provision environments come with a preloaded dataset of users. We'll use this dataset to generate a set of tickets/match making requests, and then we'll match them using Redis and RediSearch.

For this purpose we will need a number of things on Redis side: an input Stream, which players/clients would use to put their requests/tickets on. A MatchMaker, that listens to the stream and looks for potential matches. And an output Stream, for the matched games, to be given back to the client or downstream for further processing (e.g. looking for an available game server if it's not P2P, etc., etc.) We'll also emit metrics on another output Stream.

:insert diagram here:

###
First, run the program that will maintain the stats stream:
```
go run stats.go
```
Second (in a separate tab/window), run the matchmaker:

## Next steps

Well done, you made it through the first exercise! Take a short break if you want, and then move on to [exercise 2](exercise-2-start.md).
