package api

import (
	"github.com/aws/aws-sdk-go/aws/session"

	// TODO: use it
	_ "github.com/aws/aws-sdk-go/service/dynamodb"
)

// NoSQLStorage is an interface for interacting with a NoSQL
// database
type NoSQLStorage interface {
}

func testDynamoDB() {
	mySession := session.Must(session.NewSession())
	_ = mySession
}
