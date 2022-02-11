package Oris_Log

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

// SqlWriter
type SqlWriter struct {
	config *config
	db  *sql.DB
	context map[string]interface{}
	ID uuid.UUID
}

// NewContext This function is used to add context to a log record.
func (s *SqlWriter) NewContext(ctx map[string]interface{}) Logger {
	return &SqlWriter{config: s.config, context: ctx, ID: uuid.New()}
}

// AddContext Add a new context value to log writer
func (l *SqlWriter) AddContext(key string, value interface{}) {
	l.context[key] = value
}

// GetContext returns context value based on its key
func (l *SqlWriter) GetContext(key string) interface{} {
	return l.context[key]
}

// Info log info
func (s *SqlWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := s.config.Disabled["info"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, INFO, s.ID)
		go s.write(&lf)
	}
}

// Error log error
func (s *SqlWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := s.config.Disabled["error"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, ERROR, s.ID)
		go s.write(&lf)
	}
}

// Debug log debug
func (s *SqlWriter) Debug(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := s.config.Disabled["debug"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, DEBUG, s.ID)
		go s.write(&lf)
	}
}

// Fatal log fatal
func (s *SqlWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := s.config.Disabled["fatal"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, FATAL, s.ID)
		go s.write(&lf)
		time.Sleep(2e9)
		os.Exit(0)
	}
}

func (s *SqlWriter) write(lf *logFormat)  {
	ctx:=context.Background()
	context, _:=json.Marshal(lf.Context)
	insertQuery := fmt.Sprintf(`Insert into %s (Created, ID, Type, Prefix, Source, Message, Context) 
		values (?,?,?,?,?,?,?)`, s.config.DBName)
	_, err := s.db.ExecContext(ctx, insertQuery, lf.Created, lf.ID, lf.Label, lf.Prefix, lf.Source, lf.Message, string(context))
	if err != nil {
		log.Println(err)
	}
}
