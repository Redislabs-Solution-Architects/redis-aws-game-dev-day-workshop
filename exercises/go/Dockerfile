# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY /stats/stats.go ./

RUN go build -o /docker-stats

CMD [ "/docker-stats" ]