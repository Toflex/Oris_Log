package Oris_Log

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func ExampleNew() {
	log:= New()
	log.Info("Hello World!")
	// Output: {"created":"2022-01-04 01:24:19","id":"734eb0ec-37e0-4c22-bc33-6cb8f83a6734","label":"INFO","prefix":"awesomeProject","source":"main.go:17","message":"Log without context"}
}

func ExampleNew_mongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
	if err != nil {
		panic(err)
	}

	conn := client.Database("baz")

	log:= New(conn)
	log.Info("Hello World!")
}

func ExampleLogWriter_SetContext() {
	log:= New()
	ctx:= map[string]interface{}{"userId":"b60dbcda-54dc-4cc0-b46f-56457f52314b"}
	l:=log.SetContext(ctx) // Adds userId to all log
	l.Info("Hello World!")
}
