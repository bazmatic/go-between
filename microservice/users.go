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

		requestID, err := unmarshalRequest(m, userNewRequest)
		if err != nil {
			panic(err)
		}

		userNewResponse := pb.UserNewResponse{
			Error: errMessage,
			Data: &pb.User{
				Name: userNewRequest.Name,
				Id:   33,
			},
		}

		marshalledMessage, resSubject, err := marshalResponse(&userNewResponse, requestID, subjects.SubjectUserCreateCompleted)
		if err != nil {
			panic(err)
		}

		natsClient.Publish(resSubject, marshalledMessage)
		logMessageSent(resSubject)

	}); err != nil {
		log.Fatal(err)
	}

	// Wait forever
	select {}
}

//Unmarshal the request and return the attached requestID
func unmarshalRequest(m *nats.Msg, pbMessage proto.Message) (requestID string, err error) {
	err = proto.Unmarshal(m.Data, pbMessage)
	requestID = getRequestID(m.Subject, subjects.SubjectUserCreate)
	return
}

func marshalResponse(pbMessage proto.Message, requestID string, parentSubject string) (marshalledMessage []byte, subject string, err error) {
	marshalledMessage, err = proto.Marshal(pbMessage)
	subject = requestID + "." + subjects.SubjectUserCreateCompleted
	return
}

func getRequestID(subject string, parentSubject string) (requestID string) {
	// Get request ID from subject
	parts := strings.Split(subject, "."+parentSubject)
	if len(parts) != 2 {
		panic("Invalid subject")
	}
	requestID = parts[0]
	return
}

func logMessageReceived(m *nats.Msg) {
	log.Printf("MSG IN: %s", m.Subject)
}

func logMessageSent(subject string) {
	log.Printf("MSG OUT %s", subject)
}
