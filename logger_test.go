package Oris_Log

import (
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)


func Test_Log_File_Creation(t *testing.T)  {
	lg:= New(Configuration{Name: "test_log_file", Output: FILE})
	lg.Info("Hello World!")

	time.Sleep(1e9)

	filename:="sample.json"
	_, inputError := os.Open(filename)
	if inputError != nil {
		t.Fatalf("File was not created for logging. %s", inputError)
	}
}

func TestConsoleWriter_AddContext(t *testing.T) {
	lg:=New(Configuration{Name: "test_log_console"})
	traceId:=uuid.New().String()
	lg.Debug("Hello, World!")
	lg.Debug("Hello, World!")
	lg.Debug("Hello, World!")

	l:=lg.NewContext()
	l.AddContext("traceId", traceId)
	if l.GetContext("traceId") != traceId{
		t.Errorf("Context was not set with traceId=%s", traceId)
	}
	l.Info("Log with trace ID")
	ttt(l)

}

func ttt(x Logger) {
	x.AddContext("Name", "Log test")
	x.Info("testing testing!")
	x.SetLogID("15151515151515")
	x.Debug("Log with trace ID same as caller")
}
