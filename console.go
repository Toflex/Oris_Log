package Oris_Log

import (
	"encoding/json"
	"fmt"
	"os"
)

// LogWriter
type LogWriter struct {
	config *config
	context map[string]interface{}
}

func writer(log *logFormat) {
	js,_:=json.Marshal(log)
	fmt.Fprintf(os.Stdout, "%s\n", string(js))
}

// SetContext This function is used to add context to the logger instance
func (l *LogWriter) SetContext(ctx map[string]interface{}) Logger {
	return &LogWriter{config: l.config, context: ctx}
}

// Info log info
func (l *LogWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, INFO)
	writer(&lf)
}

// Error log error
func (l *LogWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, ERROR)
	writer(&lf)
}

// Warning log warning
func (l *LogWriter) Warning(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, WARNING)
	writer(&lf)
}

// Fatal log fatal
func (l *LogWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), l.config.Name, l.context, FATAL)
	writer(&lf)
	os.Exit(0)
}
