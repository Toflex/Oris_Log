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
// Oris logger requires a configuration file `logs.json`:
//  {
//  "name": "awesomeProject",
//  "filename": "sample",
//  "MaxFileSize": "2K",
//  "folder": "logs",
//  "output": "console",
//  "buffer": 10000
//	}
// 1. name: The project name
//
// 2. filename: The name that would be assigned to log file. This config is only needed when using file as output.
//
// 3. MaxFileSize: Set the max file size for each logs, ones the file exceeds the max size, it is renamed and subsequent logs are written to a new file.
// Note: A new log file is created every day. The size is a string with suffix 'K' kilobyte and 'M' Megabyte.
//
// 4. folder: the name of the folder where logs file would be kept, the default location is ./logs.
//
// 5. output: This determines where logs would be written to i.e file, console, or MongoDB
//
// 6. buffer: This id used when output is set to `file`, it buffers file during write to memory to avoid blocking.
//
package Oris_Log

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
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
	
	// Warning Logs are written to an output with a warning label
	Warning(interface{}, ...interface{})

	// Fatal : Logs is been written to output and program terminates immediately  
	Fatal(interface{}, ...interface{})

	// SetContext This function is used to add context to a log record.
	SetContext(map[string]interface{}) Logger
}

type Label string

const (
	configFileName = "logs.json"
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
	WARNING Label = "WARNING"
	FATAL   Label = "FATAL"
)

type config struct {
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	MaxFileSize string `json:"MaxFileSize"`
	MaxSize     int64  `json:"max_size"`
	HasMaxSize  bool   `json:"HasMaxSize"`
	Folder      string `json:"folder"`
	Path        string `json:"path"`
	Output      string `json:"output"`
	BaseDir     string `json:"base_dir"`
	Buf         int    `json:"buffer"` // Buf only works with file output type
	DBName      string `json:"DBName"`
}

// LogFormat outlines the way messages will be formatted in json
type logFormat struct {
	Created string    `json:"created"`
	ID     uuid.UUID `json:"id"`
	Label  Label     `json:"label"`
	Prefix string    `json:"prefix"`
	Source  string    `json:"source"`
	Message string    `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// initConfig Initializes log configurations that can be gotten from the `logs.json` file.
func (c *config) initConfig() {
	_, b, _, _ := runtime.Caller(2)
	c.BaseDir = filepath.Dir(b)

	inputFile, inputError := os.Open(fmt.Sprintf("%s/%s", c.BaseDir, configFileName))
	if inputError != nil {
		panic("Log configuration file not found")
	}
	defer inputFile.Close()

	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		panic(err)
	}

	// Get configurations file contents
	json.Unmarshal(file, c)

	if c.Filename == "" {
		c.Filename = "log.json"
	} else if !strings.HasSuffix(strings.TrimSpace(c.Filename), ".json") {
		c.Filename = fmt.Sprintf("%s.json", c.Filename)
	}
	if c.Output == "" {
		c.Output = CONSOLE
	}
	if c.DBName == ""  {
		c.DBName = DBNAME
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
}

// New create a logger instance
func New(conn ...interface{}) Logger {
	// Initialize configuration file
	configs := &config{}
	configs.initConfig()

	switch configs.Output {
	case CONSOLE:
		return &LogWriter{config: configs}
	case FILE:
		{
			file := &FileWriter{config: configs, ch: make(chan logFormat, configs.Buf)}
			go processor(file.ch, file)
			return file
		}
	case SQL:
		{
			//if len(conn) == 0 {
			//	panic("Failed to connect to SQL DB")
			//}
			//if reflect.TypeOf(conn).String() != "*sql.DB" {
			//	panic("SQL requires *sql.DB connection type")
			//}
			//db:=conn[0].(*sql.DB)
			//if db == nil{
			//	log.Fatal("Sql DB connection is nil")
			//}
			//return &SqlWriter{
			//	config: configs,
			//	db: db,
			//}
			panic("SQL writer no implemented")
		}
	case MONGODB:
		{
			if len(conn) == 0 {
				panic("Failed to connect to Mongo DB")
			}
			if reflect.TypeOf(conn).String() != "*mongo.Database" {
				panic("MongoDb requires *mongo.Collection connection type")
			}
			db:= conn[0].(*mongo.Database)
			if db == nil{
				log.Fatal("MongoDB connection is nil")
			}

			// Instance of mongo DB Collection
			col := db.Collection(configs.DBName)

			return &MongoWriter{
				config: configs,
				db: col,
			}
		}
	}

	panic("Invalid Output!")
}

//	getSource returns the name of the function caller and the line where the call was made
func getSource(depth int) string {
	_, file, line, _ := runtime.Caller(depth) // 2
	caller := filepath.Base(file)
	return fmt.Sprintf("%s:%d", caller, line)
}

func prepareLog(message, prefix string, ctx map[string]interface{}, ltype Label) logFormat {
	return logFormat{
		ID:      uuid.New(),
		Created: time.Now().Format(LOGDATEFORMAT),
		Label:   ltype,
		Prefix:  prefix,
		Source:  getSource(3),
		Context: ctx,
		Message: message,
	}
}