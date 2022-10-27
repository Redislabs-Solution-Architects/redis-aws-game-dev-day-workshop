<img src="img/redis-logo-full-color-rgb.png" height=100/><img align="right" src="img/aws-logo-1.jpeg" height=100 />

# Matchmaking your way into the leaderboards using Redis Enterprise on AWS
Hello and thank you for joining this workshop/hands-on lab! We're happy to have you join, so welcome! We hope you will enjoy this workshop and hopefully learn something new. This README will explain how to get started and point you to the exercises and how to perform them. In any case, the instructors will be here to guide you and help you out if you need any assistance or have any questions. Feedback is always appreciated! If you would like more info on the topic of matchmaking feel free to check out Volkan Civelek's excellent (and free!) [e-book on Matchmaking](https://redis.com/docs/the-game-developers-guide-to-matchmaking/) as it contains many interesting topics and best practices around this topic.

## Prerequisites

1. A modern development laptop with MacOS or Linux. Linux on Windows via WSL2 will probably work too.
1. A working internet connection
1. Docker Desktop, core components or some other way to build and run containers
1. An active AWS account (we'll be using the free Redis Enterprise Cloud tier or an instance will be provided)

## Required knowledge
Some development experience is required to go through the exercises, but it should be accessible enough for developers of any level as well as other interested folks who might not be developers. Feel free to raise any challenges to us, either directly, or by submitting an issue on this repo.

## Getting started
We have provided you with a Redis Enterprise Cloud database endpoint and the needed credentials via e-mail. You can use these to connect to a pre-provisioned instance that has all the features and data that is required for these exercises. We also provided a link to a running Grafana instance that is connected to the REC database instance and is pre-populated with the dashboards used during the exercises.

## Running locally
If you want to run things locally instead of in the cloud, you can run all the serverside components locally by running `docker compose up` from the [exercises/grafana](exercises/grafana) folder. This will run all the images and software used locally.

## Exercises

* Exercise 1 - High quality matchmaking with large volumes of players using RediSearch: [start](exercises/exercise-1-start.md)
* Exercise 2 - A real life matchmaking application using Go and Redis: [start](exercises/exercise-2-start.md)
* Exercise 3 - Fast and efficient Leaderboards using Redis: [start](exercises/exercise-3-start.md)

## Troubleshooting
Reach out to one of the instructors in case of any issues, we're always glad to help. If you're using this repo outside of a workshop or hands-on lab, feel free to submit an issue.

# License
Redis proprietary, subject to the Redis Enterprise Software and/or [Cloud Services license](https://redis.com/legal/cloud-tos/).