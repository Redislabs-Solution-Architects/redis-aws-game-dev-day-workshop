<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 2 - Fast and efficient Leaderboards using Redis
Hope you enjoyed and learned something from the exercises around matchmaking. This exercise will introduce you to the data type that powers many Leaderboards across the world. As with matchmaking, collecting stats/metrics from individual players can be a daunting task, especially with very large numbers of players. How do you keep this up-to-date in real-time and fast? How do you avoid contention on a specific metric?

There's a special data type in Redis called a Sorted Set that we can use to power a leaderboard (or any type of ranked data). Let's get started! 

## Goals

* Learn about the [Sorted Set](https://redis.io/docs/data-types/sorted-sets/) data type and its abilities
* Learn how to build a Leaderboard using a [Sorted Set](https://redis.io/docs/data-types/sorted-sets/)

## Structure of this exercise
1. Setup
1. Leaderboard basics and how to use a Sorted Set to create a Leaderboard

### Setup
As with Exercise 1, we will be connecting to the provided Redis Enterprise database and will interact with it via the Redis CLI.
### Sorted Sets
A Sorted Set (as the name implies) is an ordered collection of unique values. In Redis each value will have a score associated with it, and by updating the score as we go along the Set will maintain its ordering according to the score. Think of scenarios like maintaining a high score leaderboard when playing a game, a list of 'biggest spenders' on your bank account or other scenarios where you need to update a ranking/score as more data becomes available in your application. Especially in gaming, where there's often a very large number of players all generating information that's relevant to the leaderboard you will need a solution that's both able to show the latest updates in real time as wekk as being able to handle large amounts of updates. First we will take a look at how we can use a Sorted Set and further on in the exercise we'll also show how to get large amounts of updates from different sources in there.

* We can add members to a Set directly by using the `zadd` command. There is no need to set a key first.

The syntax of the command is `zadd <key> <score> <member>`. There are more options for this command but we'll start simple. For a full overview of all the options of `zadd` check the [documentation](https://redis.io/commands/zadd/). We will add three members to a Sorted Set using the following commands:

```
zadd lb 2 "Lars"
zadd lb 3 "Paul"
zadd lb 4 "Patricia"
zadd lb 1 "Mary"
```
* Let's get the top 3 from this Sorted Set using the following command:
```
zrange lb 0 2
```
* Notice that this the bottom 3 in the Sorted Set by their respective scores, which in this case is Mary, Lars and Paul. This doesn't really make sense from a leaderboard perspective so we can also do the reverse and get the top 3 in the Sorted Set by using:
```
zrange lb 0 2 rev
```

Of course real life is more dynamic than this, so scores will change over time. For this purpose we can use the `zincrby <key> <increment> <member>` command. Note that like `zadd` the command can be issues without the Set or any of its members existing in the first place, so in your code you don't need to make that distinction either and you can implement only the `zincrby`. Try the following command:

```
zincrby lb 321 "Mary"
```

And get the top 3 out again using:
```
zrange lb 0 2 rev
```
Notice how Mary is now at the top of the leaderboard. So from client perspective, all you have to do is submit the score for a player to your local Redis Stream and have the server component(s) read the Stream and populate the Sorted Set with the additional score. There's no need to implement any sorting or parsing logic on your application side, this will 'just work' with the out of the box Sorted Set data type in Redis. To then show the Leaderboard in a dashboard, the only thing that needs to be done is to issue a `zrange` command. To make it even easier, you can also add the `withscores` option and the command will return the top 3 with scores included. Try this out using:

```
zrange lb 0 2 rev withscores
```

It's also possible to retrieve scores/rankings close to the players own ranking/score, e.g. if our own score is 325, and we want to see which players are close to us in score, we can do so using:

```
zrange lb 350 300 rev byscore withscores
```

This will return all members with a score between 350-300. Note that because of the `rev` option the lower/upper boundary are also in reverse; if you omit `rev` you need to reverse these.

For more information on Sorted Sets and their assocatied commands, check the [Redis documentation](https://redis.io/commands#sorted_set).

Now that you have seen a few of the most commonly used Redis data structures in action, feel free to take a look at what other data types and command are available in the [Redis documentation](https://redis.io/topics/data-types)

## Next steps

How can we build Redis-powered applications and services that use sorted sets for our leaderboards?

In the final exercise we'll pull together our knowledge of RediSearch and Sorted Sets. We'll see some code that's used for matchmaking and leaderboards, introduce Redis Streams and visualise everything in Grafana.

Click [here](exercise-3-start.md) to move onto exercise 3.