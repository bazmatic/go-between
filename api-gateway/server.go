package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/labstack/echo"
	"github.com/nats-io/nats"

	pb "../proto"
	subjects "../protocol"
)

func main() {

	const APIPrefix = "/v1/api/"

	// Start NATS connection
	// TODO: Centralise NATS config
	natsClient, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to NATS")
	}
	defer natsClient.Close()

	//=== Configure endpoints
	e := echo.New()

	// POST users
	e.POST(APIPrefix+"users", func(context echo.Context) error {
		defer func() {
			if recoveredErr := recover(); recoveredErr != nil {
				echo.NewHTTPError(http.StatusInternalServerError, recoveredErr)
			}
		}()
		return reqUserPost(context, e, natsClient)
	})

	// GET users
	e.GET(APIPrefix+"users", func(context echo.Context) error {
		defer func() {
			if recoveredErr := recover(); recoveredErr != nil {
				echo.NewHTTPError(http.StatusInternalServerError, recoveredErr)
			}
		}()
		return errors.New("Not implemented")
	})

	//=== Start web server
	e.Logger.Fatal(e.Start(":3000"))

}

func reqUserPost(context echo.Context, e *echo.Echo, natsClient *nats.Conn) error {

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

		// Use the response
		var userNewResponse *pb.UserNewResponse = new(pb.UserNewResponse)
		err = proto.Unmarshal(msg.Data, userNewResponse)

		return context.JSON(http.StatusOK, userNewResponse)
	}
}

// Send request to NATS, await a response from microservice, return it
func awaitRequest(reqSubject string, resSubject string, marshalledMessage []byte, natsClient *nats.Conn) nats.Msg {
	requestID := strconv.Itoa(rand.Intn(100000))
	uniqueReqSubject := requestID + "." + reqSubject
	uniqueResSubject := requestID + "." + resSubject
	log.Printf("Req subject %s", uniqueReqSubject)
	log.Printf("Res subject %s", uniqueResSubject)
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
