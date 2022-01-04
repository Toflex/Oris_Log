package Oris_Log

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

// MongoWriter
type MongoWriter struct {
	config  *config
	db      *mongo.Collection
	context map[string]interface{}
}

// SetContext This function is used to add context to a log record.
func (m *MongoWriter) SetContext(ctx map[string]interface{}) Logger {
	return &MongoWriter{config: m.config, context: ctx}
}

// Info log info
func (m *MongoWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, INFO)
	m.writer(&lf)
}

// Error log error
func (m *MongoWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, ERROR)
	writer(&lf)
}

// Warning log warning
func (m *MongoWriter) Warning(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, WARNING)
	writer(&lf)
}

// Fatal log fatal
func (m *MongoWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, FATAL)
	writer(&lf)
	os.Exit(0)
}

func (m *MongoWriter) writer(l *logFormat) {
	_, err := m.db.InsertOne(context.Background(), &l)
	if err != nil {
		log.Println(err)
	}
}
