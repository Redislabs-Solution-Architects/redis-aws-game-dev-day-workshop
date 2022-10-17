<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 1 - Leaderboards
Hello and thank you for joining this hands-on lab! We're happy to have you join, so welcome! This first exercise aims to get you up and running and introduce you to some of the core concepts and show you how to build and query a Leaderboard. If you're stuck, don't dwell too long on it. Instead, check out the hints and the solution [over here](exercise-1-solution.md) or reach out to one of the instructors, who will be quite happy to help!

## Goals

* Learn some of the basic Redis commands by using the [Redis CLI](https://redis.io/topics/rediscli)
* Learn about some of the different [data types in Redis](https://redis.io/topics/data-types-intro)
* Learn about the [Hash](https://redis.io/docs/data-types/hashes/) data type and how we will use it to manage Tickets
* Learn how to build a Leaderboard using a [Sorted Set](https://redis.io/docs/data-types/sorted-sets/)
* Learn how to use [RediSearch](https://redis.io/docs/stack/search/) to execute fast queries on Redis Hashes or [JSON](https://redis.io/docs/stack/json/) documents

### Starting Redis and opening the CLI
* In a terminal Window, open the Redis CLI and connect to your provided Redis Enterprise database by typing:
```
redis-cli -h <hostname> -p <port>
```

By default, the Redis CLI will connect to 127.0.0.1 and port 6379. If you want to use a different hostname and port, you can use the `-h` and `-p` options to specify a diferent hostname and port. Use the `--tls` and `--sni <hostname>` options in case of TLS and use the `--insecure` option when self-signed certificates are used. Those are not needed for the provided environments however.

### Strings and hashes
* Let's see if we can store something in Redis and retrieve it again, so let's start simple and start with a basic String key/value pair.
* Add a String key/value pair to Redis with a key of `hello` and a value of `aws`:
```
set hello aws
```
* Retrieve the key/value pair from Redis
```
get hello
```
Congratulations, you are now a Redis veteran! Of course, `get` and `set` are not the only two commands available for dealing with Strings in Redis. You can find out about all of the other commands for Strings (and all the other datatypes) by visiting the [Redis documentation](https://redis.io/commands/#string)

Let's check out some of the other data structures in Redis as well, as we will be using some of these during the exercises. Let's start with a Hash. Hashes are typically used to store flat structures with multiple attributes, such as records or structs. A very common use case is to store an HTTP or user session in Redis so you can keep your own application stateless and add/remove instances/pods as needed while still being able to retrieve state from a logged in user. It's typically much faster and more efficient to store this in a solution that's optimised for fast memory access rather than a general purpose solution.

In our exercises we will be using Hashes to store our Game tickets, which we can then query using RediSearch.
* Let's add a Hash with the key `myhash` and two key/value pairs within that Hash with keys `hello` and `how` and values `world` and `areyou?`:
```
hset myhash hello world how areyou?
```
* Retrieve the Hash and all its key/value pairs from Redis
```
hgetall myhash
```
* Retrieve a single key from the Hash
```
hget myhash hello
```
Did you notice that Redis has different commands for different data structures? For instance, a `get` command will not work on a hash, but it will work on all other data structures. And an `hget` will work on a hash but not on another data structure. If you try this you will get a 'WRONGTYPE' error. Don't know the type of a certain key? You can ask Redis the type of the `myhash` key by typing:
```
type myhash
```
and Redis will tell you that it's a hash. Alternatively, you can type:
```
type hello
```
To find out that this key is a String type. If you want to learn more about hashes then check out all the available commands at the [Redis Documentation](https://redis.io/commands#hash). You can use the drop down list on that [same page](https://redis.io/commands) to select a different data structure to find out what commands apply to that particular data structure.

### Sorted Sets
Now, let's take a look at a Sorted Set. A Sorted Set (as the name implies) is an ordered collection of unique values. In Redis each value will have a score associated with it, and by updating the score as we go along the Set will maintain its ordering according to the score. Think of scenarios like maintaining a high score leaderboard when playing a game, a list of 'biggest spenders' on your bank account or other scenarios where you need to update a ranking/score as more data becomes available in your application.

* We can add members to a Set directly by using the `zadd` command. There is no need to set a key first. So let's add three members to a Sorted Set using the following commands:
```
zadd mysortedset 2 "two"
zadd mysortedset 3 "three"
zadd mysortedset 1 "one"
```
* Notice how the order is off when adding the members to the set (we're adding the second one first, the third one second and the first one third). Let's get the first two members of the Sorted Set by typing:
```
zrange mysortedset 0 1
```
* Notice that this returns "one" and "two". We can also get the last two members of the Sorted Set by reversing the range:
```
zrange mysortedset 0 1 rev
```
This will produce the result "three" and "two". For more information on Sorted Sets and their assocatied commands, check the [Redis documentation](https://redis.io/commands#sorted_set).

Now that you have seen a few of the most commonly used Redis data structures in action, feel free to take a look at what other data types and command are available in the [Redis documentation](https://redis.io/topics/data-types)

### Modules
 But wait, there's more! While Redis has a very powerful set of data structures straight out of the box, it's always possible that you find yourself in a situation where these are not enough to cover your use cases. In those situations Redis Modules comes to the rescue! Modules allow you to extend Redis data structures and features in a modular way; you only add those modules to a specific database that you want and no more. This keeps your Redis setup fast and lean. There are many community modules available and Redis provides several modules itself as well. If that's not enough you can also write your own. The Docker container that we used to start Redis has a few modules bundled with it already, so let's take a look at a few examples. For more information on modules, see the [Redis Modules Hub](https://redis.com/community/redis-modules-hub/) and the [Redis Enterprise Modules](https://redis.com/redis-enterprise/modules/) pages.

### RediSearch
Let's start with RediSearch, a full-text search module for Redis. Retrieving keys by their primary value is fine for many use cases, but what if I have data in the cache that needs to be searched? E.g. a product model, or stores, or transactions by a certain vendor, etc. Typically this is where a key/value model starts to show its limitations. RediSearch to the rescue! RediSearch adds full-text search capabilities to Redis as well as a lot more. See the [RediSearch documentation](https://docs.redis.com/latest/modules/redisearch/) for more detailed information on this module, or checkout the [GitHub repo](https://github.com/RediSearch/RediSearch).

* First, we'll create a search index on Hash structures matching a certain prefix:
```
ft.create my_idx on hash prefix 1 my schema hello text phonetic dm:en
```
This command is a little bit more elaborate than the previous ones, so let's explore it in detail a bit more. We're creating an index called `my_idx` on the `hash` datastructure with one prefix `my` (remember that we created a Hash earlier that had the key `myhash`? also note that if you used a different key for the hash that you should also use a different prefix here) and we define the schema to be on the `hello` field (this field is present in the Hash we created earlier) which we define as a `TEXT` field and we also setup phonetic search using the `Double Metaphone for English` matcher. That's quite a handful, so don't worry if that doesn't mean a lot to you right now, let's just see what it actually does!
* Let's search our immense data set of 1 Hash and see if we can find what we want by typing a few different commands:
```
ft.search my_idx "world"
```

This will return our Hash that we created earlier.

* You can also do a wildcard search:
```
ft.search my_idx "wor*"
```
* Remember that we setup the index to be phonetic? Try it out by something that sounds like `world` but not quite:
```
ft.search my_idx "wurld"
```
This is great if your users need to search for something, but may not be able to recall exactly what it was they are searching for, or how it is spelled exactly. Another great feature is that search results can be highlighted on which part the match was found:
```
ft.search my_idx "world" highlight
```
Notice how the word `world` has no been surrounded by ```<b></b>``` tags. This is great in case we want to visually display the matching words differently in our UI. We can also change the tags by doing the following:
```
ft.search my_idx "world" highlight tags <hello> </hello>
```
Notice how the word `world` is now surrounded with the tags of our choosing.

And that's not all of the module's functionality; there's plenty more, so if you want to learn more about all the functionality of the RediSearch module, please check the [RediSearch documentation](https://oss.redis.com/redisearch/).


## Next steps

Well done, you made it through the first exercise! Take a short break if you want, and then move on to [exercise 2](exercise-2-start.md).
