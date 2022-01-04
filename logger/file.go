package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FileWriter writes log to file
type FileWriter struct {
	sync.Mutex
	config *config
	acc []logFormat
	fp *os.File
	ch chan logFormat
	context map[string]interface{}
}

// SetContext This function is used to add context to a log record.
func (f *FileWriter) SetContext(ctx map[string]interface{}) Logger {
	return &FileWriter{
		config: f.config,
		fp: f.fp,
		ch: f.ch,
		context: ctx,
		acc: f.acc,
	}
}

// Info log info
func (f *FileWriter) Info(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), f.config.Name, f.context, INFO)
	f.ch <- lf
}

// Error log error
func (f *FileWriter) Error(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), f.config.Name, f.context, ERROR)
	f.ch <- lf
}

// Warning log warning
func (f *FileWriter) Warning(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), f.config.Name, f.context, WARNING)
	f.ch <- lf
}

// Fatal log fatal
func (f *FileWriter) Fatal(arg0 interface{}, arg1 ...interface{}) {
	lf := prepareLog(fmt.Sprintf(arg0.(string), arg1...), f.config.Name, f.context, FATAL)

	go func() {
		f.ch <- lf
		time.Sleep(1e9)
		os.Exit(0)
	}()
}

// writeToFile this method handle the way log is been written to file
func (f *FileWriter) writeToFile(lf logFormat) {
	// Rotate will create the log file
	if err := f.rotate(); err != nil {
		log.Println("Unable to write to file.",err)
		panic(err)
	}
	directory:= f.config.Path

	// Check if log accumulator slice is nil
	if f.acc == nil {
		file, err := ioutil.ReadFile(directory)
		if err != nil {
			panic(err)
		}

		// Here the magic happens!
		json.Unmarshal(file, &f.acc)
	}

	outputWriter := bufio.NewWriter(f.fp)
	f.acc = append(f.acc, lf)
	data, err := json.Marshal(f.acc)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(directory, data, 0644)
	if err != nil {
		panic(err)
	}
	outputWriter.Flush()
}

// ConvertSize converts MaxFileSize value in the config file to bytes.
func ConvertSize(size string) int64 {
	const ConversionValue float64 = 1024

	if len(size) == 0{
		return MAXFILESIZE
	}

	suffix := size[len(size)-1]
	prefix := size[:len(size)-1]

	if len(prefix) == 0 {
		return MAXFILESIZE
	}

	value, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil{
		return MAXFILESIZE
	}

	switch suffix {
	case 'K':
		value *= int64(ConversionValue)
	case 'M':
		value *= int64(math.Pow(ConversionValue, 2))
	}

	if value > MAXFILESIZE {
		value = MAXFILESIZE
	}

	return value
}

// Rotate Perform the actual act of rotating and reopening file.
func (f *FileWriter) rotate() (err error) {
	if f.fp == nil {
		f.createFile()
	}

	/*
		Requirements
		A new file is created if
			1. File size does not exceed max size if HasMaxSize is true
			2. If file modified time is a that of a previous day, rename existing file if it exists and create a new file
	*/
	stat, err := f.fp.Stat()
	if err != nil {
		return err
	}

	// Conditions
	size := stat.Size()
	if f.config.HasMaxSize && size >= f.config.MaxSize {
		// Create log file
		err := f.createFile()
		if err != nil {
			 return err
		}
	} else if stat.ModTime().Format(DATEFORMAT) != time.Now().Format(DATEFORMAT) {
		// Create log file
		err := f.createFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FileWriter) createFile() error {
	f.Lock()
	defer f.Unlock()
	err:=os.MkdirAll(f.config.BaseDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	if f.fp == nil{
		// Create a file.
		fp, err := os.OpenFile(f.config.Path,os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil{
			panic(err)
		}

		file, err := ioutil.ReadFile(f.config.Path)
		if err != nil {
			panic(err)
		}

		// Here the magic happens!
		json.Unmarshal(file, &f.acc)
		f.fp = fp
		return nil
	}

	// Rename dest file if it already exists
	stat, err := os.Stat(f.config.Path)
	if err != nil {
		return nil
	}

	date:=time.Now().Format(DATEFORMAT)
	if stat.ModTime().Format(DATEFORMAT) != time.Now().Format(DATEFORMAT) {
		date= stat.ModTime().Format(DATEFORMAT)
	}

	filename := strings.TrimSuffix(f.config.Path, ".json")
	newName:=fmt.Sprintf("%s_%s.json",filename,date)

	_, exist := os.Stat(newName)
	count:=1
	for exist == nil {
		newName=fmt.Sprintf("%s_%s_(%d).json",filename,date,count)
		_, exist = os.Stat(newName)
		count++
	}
	_, err = os.Stat(f.config.Path)
	err = os.Rename(f.config.Path, newName)
	if err != nil {
		return err
	}
	f.acc = nil
	fp, err := os.OpenFile(f.config.Path,os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil{
		panic(err)
	}
	f.fp = fp

	return nil
}

func processor(ch <-chan logFormat, f *FileWriter)  {
	for {
		lf := <- ch
		f.writeToFile(lf)
	}
}