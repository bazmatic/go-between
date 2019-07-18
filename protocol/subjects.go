package protocol

// SubjectUserCreate - NATS subject string for creating a new user
const SubjectUserCreate string = "User.create"

// SubjectUserCreateCompleted - NATS subject string for confirming creation of a new user
const SubjectUserCreateCompleted string = "User.create.completed"

// SubjectUserList - NATS subject string for listing users
const SubjectUserList string = "User.list"

// SubjectUserListCompleted - NATS subject string containing a list of users
const SubjectUserListCompleted string = "User.list.completed"
