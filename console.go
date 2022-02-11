package Oris_Log

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
)

// ConsoleWriter writes log to console
type ConsoleWriter struct {
	ID uuid.UUID
	config *config
	context map[string]interface{}
}

func writer(log *logFormat) {
	js,_:=json.Marshal(log)
	fmt.Fprintf(os.Stdout, "%s\n", string(js))
}

// NewContext This function is used to add context to the logger instance
func (l *ConsoleWriter) NewContext(ctx map[string]interface{}) Logger {
	return &ConsoleWriter{config: l.config, context: ctx, ID: uuid.New()}
}

// AddContext Add a new context value to log writer
func (l *ConsoleWriter) AddContext(key string, value interface{}) {
	l.context[key] = value
}

// GetContext returns context value based on its key
func (l *ConsoleWriter) GetContext(key string) interface{} {
	return l.context[key]
}

// Info log info
func (l *ConsoleWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := l.config.Disabled["info"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, INFO, l.ID)
		writer(&lf)
	}
}

// Error log error
func (l *ConsoleWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := l.config.Disabled["error"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, ERROR, l.ID)
		writer(&lf)
	}
}

// Debug log debug
func (l *ConsoleWriter) Debug(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := l.config.Disabled["debug"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, DEBUG, l.ID)
		writer(&lf)
	}
}

// Fatal log fatal
func (l *ConsoleWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	if _,ok := l.config.Disabled["fatal"]; !ok {
		lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, FATAL, l.ID)
		writer(&lf)
		os.Exit(0)
	}
}
