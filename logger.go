// Package Oris_Log allows users to write application logs to file, console or DB (NoSql/Sql)
//
// The logger allows users to add context (context are more information that can be shared across the logger instance) to logs if needed.
//
// Features
//
//1. Write log to SQL/NoSQL DB.
//
//2. Write log to file.
//
//3. Write log to console.
//
//4. Logs are written in json/plain text format.
//
//5. Logs can be written to multiple output source i.e. to both file and console.
//
//6. Context can be added to a log i.e. User ID, can be added to a log to track user footprint.
//
// Oris logger requires a configuration file `logger.json`:
//  {
//  "name": "awesomeProject",
//  "filename": "sample",
//  "MaxFileSize": "2K",
//  "folder": "logs",
//  "output": "console",
//  "buffer": 10000,
//  "disable_logType": ["fatal","debug"]
//	}
// 1. name: The project name
//
// 2. filename: The name that would be assigned to log file. This config is only needed when using file as output.
//
// 3. MaxFileSize: Set the max file size for each log, ones the file exceeds the max size, it is renamed and subsequent logs are written to a new file.
// Note: A new log file is created every day. The size is a string with suffix 'K' kilobyte and 'M' Megabyte.
//
// 4. folder: the name of the folder where logs file would be kept, the default location is ./logs.
//
// 5. output: This determines where logs would be written to i.e file, console, or MongoDB
//
// 6. buffer: This value is used when output is set to `file`, it buffers file during write to memory to avoid blocking.
//
// 7. disable_logType: This can be used to disable so log type in production i.e. debug log. The allowable values are 'info','debug','error','fatal'
//
package Oris_Log

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Logger An interface that states the requirements for a writer to be an instance of Logger.
// Every log writer must implement the Logger interface.
type Logger interface {
	// Info Logs are written to an output with an info label
	Info(interface{}, ...interface{})

	// Error Logs are written to an output with an error label
	Error(interface{}, ...interface{})
	
	// Debug Logs are written to an output with a debug label
	Debug(interface{}, ...interface{})

	// Fatal : Logs is been written to output and program terminates immediately  
	Fatal(interface{}, ...interface{})

	// NewContext This function is used to add context to a log record, a new instance of the log writer is returned.
	NewContext() Logger

	// AddContext Add a new context value to a log writer
	AddContext(key string, value interface{})

	// GetContext returns context value based on its key
	GetContext(key string) interface{}

	// SetLogID set ID for the current log context
	SetLogID(id string)
}

// Label datatype for log type
type Label string

const (
	CONSOLE        = "console"
	FILE           = "file"
	SQL            = "sql"
	MONGODB        = "mongo"
	LOGDATEFORMAT  = "2006-01-02 15:04:05"
	DATEFORMAT     = "2006_01_02"
	DBNAME 		   = "Log"

	MAXFILESIZE int64 = 214748364800

	INFO    Label = "INFO"
	ERROR   Label = "ERROR"
	DEBUG   Label = "DEBUG"
	FATAL   Label = "FATAL"
)

type config struct {
	Name        string
	Filename    string
	MaxFileSize string
	MaxSize     int64
	HasMaxSize  bool
	Folder      string
	Path        string
	Output      string
	BaseDir     string
	Buf         int 	// Buf only works with file output type
	Collection  string
	Disabled    map[string]bool
	DisableLogType []string // values 'info','debug','error','fatal'
}

type Configuration struct {
	Name string
	Filename string
	MaxFileSize string
	Folder string
	Output string
	Buffer int
	DisableLogType []string
	MongoDB *mongo.Database
	ColName string
}

// LogFormat outlines the way messages will be formatted in json
type logFormat struct {
	Created string    `json:"created"`
	ID     string `json:"id"`
	Label  Label     `json:"label"`
	Prefix string    `json:"prefix"`
	Source  string    `json:"source"`
	Message string    `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// initConfig Initializes log configurations that can be gotten from the `logger.json` file.
func (c *config) initConfig(conf *Configuration) {
	_, b, _, _ := runtime.Caller(2)
	c.BaseDir = filepath.Dir(b)

	c.Name = conf.Name
	c.Filename = conf.Filename
	c.MaxFileSize = conf.MaxFileSize
	c.Folder = conf.Folder
	c.Output = conf.Output
	c.Buf = conf.Buffer
	c.DisableLogType = conf.DisableLogType
	c.Collection = conf.ColName

	if c.Filename == "" {
		c.Filename = "log.json"
	} else if !strings.HasSuffix(strings.TrimSpace(c.Filename), ".json") {
		c.Filename = fmt.Sprintf("%s.json", c.Filename)
	}
	if c.Output == "" {
		c.Output = CONSOLE
	}
	if c.Collection == ""  {
		c.Collection = DBNAME
	}
	if strings.TrimSpace(c.Folder) == "" {
		c.Folder = "logs"
	}

	if strings.TrimSpace(c.MaxFileSize) != "" {
		c.HasMaxSize = true
		c.MaxSize = ConvertSize(c.MaxFileSize)
	}

	c.BaseDir = fmt.Sprintf("%s/%s", c.BaseDir, c.Folder)
	c.Path = fmt.Sprintf("%s/%s", c.BaseDir, c.Filename)

	if c.Buf == 0 || c.Buf < 5000 {
		c.Buf = 10000
	}

	c.Disabled = make(map[string]bool)
	if len(c.DisableLogType) > 0 {
		for _, v := range c.DisableLogType[:4] {
			c.Disabled[v] = true
		}
	}
}

// New create a logger instance
func New(conf Configuration) Logger {
	// Initialize configuration file
	configs := &config{}
	configs.initConfig(&conf)

	switch configs.Output {
	case CONSOLE:
		return &ConsoleWriter{config: configs, ID: uuid.New().String()}
	case FILE:
		{
			//panic("File writer not implemented")
			file := &FileWriter{config: configs, ch: make(chan logFormat, configs.Buf), ID: uuid.New().String()}
			go processor(file.ch, file)
			return file
		}
	case SQL:
		{
			panic("SQL writer not implemented")
		}
	case MONGODB:
		{
			panic("Mongo writer not implemented")
			// Instance of mongo DB Collection
			col := conf.MongoDB.Collection(configs.Collection)
			return &MongoWriter{
				config: configs,
				db: col,
				ID: uuid.New().String(),
			}
		}
	}

	return &ConsoleWriter{config: configs, ID: uuid.New().String()}
}

//	getSource returns the name of the function caller and the line where the call was made
func getSource(depth int) string {
	_, file, line, _ := runtime.Caller(depth) // 2
	caller := filepath.Base(file)
	return fmt.Sprintf("%s:%d", caller, line)
}

func prepareLog(message, prefix string, ctx map[string]interface{}, ltype Label, ID string) logFormat {
	return logFormat{
		ID:      ID,
		Created: time.Now().Format(LOGDATEFORMAT),
		Label:   ltype,
		Prefix:  prefix,
		Source:  getSource(3),
		Context: ctx,
		Message: message,
	}
}