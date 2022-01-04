package Oris_Log

import (
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
