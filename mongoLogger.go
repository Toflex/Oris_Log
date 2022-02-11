package Oris_Log

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

// MongoWriter
type MongoWriter struct {
	config  *config
	db      *mongo.Collection
	context map[string]interface{}
	ID uuid.UUID
}

// NewContext This function is used to add context to a log record.
func (m *MongoWriter) NewContext(ctx map[string]interface{}) Logger {
	return &MongoWriter{config: m.config, context: ctx, ID: uuid.New()}
}

// AddContext Add a new context value to log writer
func (l *MongoWriter) AddContext(key string, value interface{}) {
	l.context[key] = value
}

// GetContext returns context value based on its key
func (l *MongoWriter) GetContext(key string) interface{} {
	return l.context[key]
}

// Info log info
func (m *MongoWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := m.config.Disabled["info"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, INFO, m.ID)
		m.writer(&lf)
	}
}

// Error log error
func (m *MongoWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := m.config.Disabled["error"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, ERROR, m.ID)
		writer(&lf)
	}
}

// Debug log debug
func (m *MongoWriter) Debug(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := m.config.Disabled["debug"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, DEBUG, m.ID)
		writer(&lf)
	}
}

// Fatal log fatal
func (m *MongoWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := m.config.Disabled["fatal"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), m.config.Name, m.context, FATAL, m.ID)
		writer(&lf)
		os.Exit(0)
	}
}

func (m *MongoWriter) writer(l *logFormat) {
	_, err := m.db.InsertOne(context.Background(), &l)
	if err != nil {
		log.Println(err)
	}
}
