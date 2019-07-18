# go-between
Experiment with Go, Echo, NATS and microservices

# Pre-requisites

- Postgres (https://www.postgresql.org/)
- Go (https://golang.org/) 
- NATS (https://nats.io/)

The following Go modules:
```
go get github.com/lib/pq
go get github.com/nats-io/nats
go get github.com/golang/protobuf/proto
go get github.com/labstack/echo

```
# Setup
## NATS
Start NATS server with default settings

## Postgres
Start Postgres
Create DB `between`:
```
CREATE DATABASE between;
```
Create table `between.users`:
```
CREATE TABLE users(
   id SERIAL,
   name VARCHAR NOT NULL
);
```

# Run
In root folder:
`make build`

## API gateway
`go run ./api-gateway/server.go`

## Microservice
`go run ./microservice/user.go`



