<img src="../img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="../img/aws-logo-1.jpeg" height=100 />

# Exercise 1b - Matchmaking in code

### How to do this all in code

Now that we've covered some of the basic queries needed for matchmaking, let's see how we can integrate that into a (simple) program. The provision environments come with a preloaded dataset of users. We'll use this dataset to generate a set of tickets/match making requests, and then we'll match them using Redis and RediSearch.

For this purpose we will need a number of things on Redis side: 

- an input Stream, which players/clients would use to put their requests/tickets on.
- A MatchMaker, that listens to the stream and looks for potential matches. 
- And an output Stream, for the matched games, to be given back to the client or downstream for further processing (e.g. looking for an available game server if it's not P2P, etc., etc.) 
- We'll also emit metrics on another output Stream.

![Diagram of matchmaker](/img/diagram.png)

###
First, run the program that will maintain the stats stream:
```
go run stats.go
```
Second (in a separate tab/window), run the matchmaker: