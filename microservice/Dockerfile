# Use an official Python runtime as a parent image
FROM golang:latest
COPY . /go/src/github.com/bazmatic/go-between/
WORKDIR /go/src/github.com/bazmatic/go-between/

RUN ["/bin/bash", "-c", "go get github.com/lib/pq"]
RUN ["/bin/bash", "-c", "go get github.com/nats-io/nats"]
RUN ["/bin/bash", "-c", "go get github.com/golang/protobuf/proto"]

WORKDIR /go/src/github.com/bazmatic/go-between/microservice

CMD ./start.sh
