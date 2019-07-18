FROM golang:latest
COPY . /go/src/github.com/bazmatic/go-between/
WORKDIR /go/src/github.com/bazmatic/go-between/
RUN ["/bin/bash", "-c", "go get github.com/lib/pq"]
RUN ["/bin/bash", "-c", "go get github.com/nats-io/nats"]
RUN ["/bin/bash", "-c", "go get github.com/golang/protobuf/proto"]
RUN ["/bin/bash", "-c", "go get github.com/labstack/echo"]
RUN ["/bin/bash", "-c", "go get github.com/dgrijalva/jwt-go"]

WORKDIR /go/src/github.com/bazmatic/go-between/api-gateway

EXPOSE 3000

CMD ./start.sh
#CMD ["go", "run", "./server.go"]