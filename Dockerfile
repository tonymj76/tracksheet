# go setup
FROM golang:alpine AS base

WORKDIR /app

COPY . .

RUN cd crawler-service && go get -d -v && go build -o crawler


# app setup
FROM alpine

WORKDIR /app
COPY --from=base /app/crawler-service/crawler /app/
COPY --from=base /app/crawler-service/env.json /app/

CMD [ "./crawler" ]