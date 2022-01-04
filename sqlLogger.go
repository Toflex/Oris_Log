package Oris_Log

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// SqlWriter
type SqlWriter struct {
	config *config
	db  *sql.DB
	context map[string]interface{}
}

// SetContext This function is used to add context to a log record.
func (s *SqlWriter) SetContext(ctx map[string]interface{}) Logger {
	return &SqlWriter{config: s.config, context: ctx}
}

// Info log info
func (s *SqlWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, INFO)
	go s.write(&lf)
}

// Error log error
func (s *SqlWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, ERROR)
	go s.write(&lf)
}

// Warning log warning
func (s *SqlWriter) Warning(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, WARNING)
	go s.write(&lf)
}

// Fatal log fatal
func (s *SqlWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), s.config.Name, s.context, FATAL)
	go s.write(&lf)
	time.Sleep(2e9)
	os.Exit(0)
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
