package main

import (
	//"errors"
	pb "../proto"
	subjects "../protocol"
	proto "github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats"
	"log"
	"strings"
)

func main() {

	// Start NATS connection
	// TODO: Centralise NATS config
	natsClient, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Microservice started")
	defer natsClient.Close()

	//=== Subscriptions
	if _, err := natsClient.Subscribe("*."+subjects.SubjectUserCreate, func(m *nats.Msg) {
		logMessageReceived(m)
		var userNewRequest *pb.UserNewRequest = new(pb.UserNewRequest)
		var errMessage string
		err = proto.Unmarshal(m.Data, userNewRequest)
		if err != nil {
			errMessage = err.Error()
		}

		// Get request ID
		parts := strings.Split(m.Subject, "."+subjects.SubjectUserCreate)
		if len(parts) != 2 {
			panic("Invalid subject")
		}
		requestID := parts[0]

		userNewResponse := pb.UserNewResponse{
			Error: errMessage,
			Data: &pb.User{
				Name: userNewRequest.Name,
				Id:   33,
			},
		}
		marshalledMessage, _ := proto.Marshal(&userNewResponse)
		resSubject := requestID + "." + subjects.SubjectUserCreateCompleted
		natsClient.Publish(resSubject, marshalledMessage)
		logMessageSent(resSubject)

	}); err != nil {
		log.Fatal(err)
	}

	// Wait forever
	select {}
}

func logMessageReceived(m *nats.Msg) {
	log.Printf("MSG IN: %s", m.Subject)
}

func logMessageSent(subject string) {
	log.Printf("MSG OUT %s", subject)
}
