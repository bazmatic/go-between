package main

import (
	"github.com/labstack/echo/middleware"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"github.com/nats-io/nats"

	pb "github.com/bazmatic/go-between/proto"
	subjects "github.com/bazmatic/go-between/protocol"
)

func main() {

	const APIPrefix = "/v1/api/"

	// Start NATS connection
	// TODO: Centralise NATS config
	//natsClient, err := nats.Connect(nats.DefaultURL)
	natsClient, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to NATS")
	}
	defer natsClient.Close()

	//=== Configure endpoints
	e := echo.New()

	// POST users
	e.POST(APIPrefix+"users", func(context echo.Context) (result error) {
		return reqUsersPost(context, e, natsClient)
	})

	// GET users
	e.GET(APIPrefix+"users", func(context echo.Context) error {
		return reqUsersGet(context, e, natsClient)
	})

	// Ensure all errors are handled as an HTTP response
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = customHTTPErrorHandler

	//=== Start web server
	e.Logger.Fatal(e.Start(":3000"))
}

// Standard error response
// {
//    Error: "<Error message>"
// }
func customHTTPErrorHandler(err error, c echo.Context) {
	c.JSON(http.StatusInternalServerError, struct{ Error string }{err.Error()})
}

func reqUsersPost(context echo.Context, e *echo.Echo, natsClient *nats.Conn) error {

	userNew := new(pb.UserNewRequest)
	err := context.Bind(userNew)
	if err != nil {
		panic(err)
	} else {
		marshalledMessage, err := proto.Marshal(userNew)
		if err != nil {
			panic(err)
		}

		msg := awaitRequest(subjects.SubjectUserCreate, subjects.SubjectUserCreateCompleted, marshalledMessage, natsClient)
		var userNewResponse = new(pb.UserNewResponse)
		err = proto.Unmarshal(msg.Data, userNewResponse)

		return context.JSON(http.StatusOK, userNewResponse)
	}
}

func reqUsersGet(context echo.Context, e *echo.Echo, natsClient *nats.Conn) error {

	usersAll := new(pb.UsersAllRequest)
	err := context.Bind(usersAll)
	if err != nil {
		panic(err)
	} else {
		marshalledMessage, err := proto.Marshal(usersAll)
		if err != nil {
			panic(err)
		}

		msg := awaitRequest(subjects.SubjectUserList, subjects.SubjectUserListCompleted, marshalledMessage, natsClient)

		// Use the response
		var usersAllResponse = new(pb.UsersAllResponse)
		err = proto.Unmarshal(msg.Data, usersAllResponse)

		return context.JSON(http.StatusOK, usersAllResponse)
	}
}

// Send request to NATS, await a response from microservice, return it
func awaitRequest(reqSubject string, resSubject string, marshalledMessage []byte, natsClient *nats.Conn) nats.Msg {

	requestID := strconv.Itoa(rand.Intn(100000))
	uniqueReqSubject := requestID + "." + reqSubject
	uniqueResSubject := requestID + "." + resSubject
	natsClient.Publish(uniqueReqSubject, marshalledMessage)
	logMessageSent(uniqueReqSubject)
	subcriber, err := natsClient.SubscribeSync(uniqueResSubject)

	if err != nil {
		log.Printf("awaitRequest: Failed to subscribe: %s", err)
		panic(err)
	}
	natsClient.Flush()
	log.Printf("Sent message to NATS")

	// Wait for response
	msg, err := subcriber.NextMsg(10 * time.Second)
	if err != nil {
		log.Printf("awaitRequest: Failed to get response: %s", err)
		panic(err)
	}
	logMessageReceived(msg)
	return *msg
}

func logMessageReceived(m *nats.Msg) {
	log.Printf("MSG IN: %s", m.Subject)
}

func logMessageSent(subject string) {
	log.Printf("MSG OUT %s", subject)
}
