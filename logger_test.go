package Oris_Log

import (
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)


func Test_Log_File_Creation(t *testing.T)  {
	lg:= New()
	lg.Info("Hello World!")

	time.Sleep(1e9)

	filename:="sample.json"
	_, inputError := os.Open(filename)
	if inputError != nil {
		t.Fatalf("File was not created for logging")
	}
}

func TestConsoleWriter_AddContext(t *testing.T) {
	lg:=New()
	traceId:=uuid.New().String()
	lg.Debug("Hello, World!")
	lg.Debug("Hello, World!")
	lg.Debug("Hello, World!")
	ctx:=make(map[string]interface{})
	l:=lg.NewContext(ctx)
	l.AddContext("traceId", traceId)
	if l.GetContext("traceId") != traceId{
		t.Errorf("Context was not set with traceId=%s", traceId)
	}
	l.Info("Log with trace ID")
	ttt(l)

}

func ttt(x Logger) {
	x.AddContext("Name", "Tolu")
	x.Info("testing testing!")
	x.SetLogID("15151515151515")
	x.Debug("Log with trace ID same as caller")
}
