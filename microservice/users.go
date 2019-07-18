package main

import (
	//"errors"
	"database/sql"
	pb "github.com/bazmatic/go-between/proto"
	subjects "github.com/bazmatic/go-between/protocol"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats"
	"log"
	"strings"
)

func main() {

	// Connect to NATS
	// TODO: Centralise NATS config
	natsClient, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer natsClient.Close()
	log.Printf("Connected to NAT")

	// Connect to DB
	connStr := "dbname=between sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to DB")

	//=== Subscriptions
	//user.create
	if _, err := natsClient.Subscribe("*."+subjects.SubjectUserCreate, func(m *nats.Msg) {

		logMessageReceived(m)
		var userNewRequest = new(pb.UserNewRequest)

		requestID, err := unmarshalRequest(m, subjects.SubjectUserCreate, userNewRequest)
		if err != nil {
			panic(err)
		}

		//Save to DB
		var userID int32
		var sql = "insert into users (name) values('" + userNewRequest.Name + "') RETURNING id"
		err = db.QueryRow(sql).Scan(&userID)

		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		userNewResponse := pb.UserNewResponse{
			Data: &pb.User{
				Name: userNewRequest.Name,
				Id:   userID,
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

	//user.list
	if _, err := natsClient.Subscribe("*."+subjects.SubjectUserList, func(m *nats.Msg) {

		logMessageReceived(m)
		var usersAllRequest = new(pb.UsersAllRequest)

		requestID, err := unmarshalRequest(m, subjects.SubjectUserList, usersAllRequest)
		if err != nil {
			panic(err)
		}

		rows, err := db.Query("select * from users")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		userAllResponse := pb.UsersAllResponse{
			//Error: err,
			Data: make([]*pb.User, 0),
		}

		for rows.Next() {
			var (
				id   int32
				name string
			)
			if err := rows.Scan(&id, &name); err != nil {
				log.Fatal(err)
			}
			userAllResponse.Data = append(
				userAllResponse.Data,
				&pb.User{
					Name: name,
					Id:   id,
				},
			)
		}
		marshalledMessage, resSubject, err := marshalResponse(&userAllResponse, requestID, subjects.SubjectUserListCompleted)
		if err != nil {
			panic(err)
		}
		natsClient.Publish(resSubject, marshalledMessage)
		logMessageSent(resSubject)

	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("Microservice ready")

	// Wait forever
	select {}
}

//Unmarshal the request and return the attached requestID
func unmarshalRequest(m *nats.Msg, parentSubject string, pbMessage proto.Message) (requestID string, err error) {
	err = proto.Unmarshal(m.Data, pbMessage)
	requestID = getRequestID(m.Subject, parentSubject)
	return
}

func marshalResponse(pbMessage proto.Message, requestID string, parentSubject string) (marshalledMessage []byte, subject string, err error) {
	marshalledMessage, err = proto.Marshal(pbMessage)
	subject = requestID + "." + parentSubject
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
	log.Printf("MSG IN  %s", m.Subject)
}

func logMessageSent(subject string) {
	log.Printf("MSG OUT %s", subject)
}
