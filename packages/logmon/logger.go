package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
)

type Logger struct{
	logger *log.Logger
	Level LogLevel
	Prefix string
}

func NewLogger(level LogLevel, prefix string) *Logger {
	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 06666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger := log.New(logFile, prefix, log.LstdFlags)
	return &Logger{
		logger: logger,
		Level: level,
		Prefix: prefix,
	}
	
}

func (l* Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}){
	if l.Level <= INFO {
		l.logger.Println("[INFO]", fmt.Sprintf("%s %v", msg, formatKeysAndValues(keysAndValues...)))
	}
}
func(l *Logger) Error(msg string, keysAndValues ...interface{}){
	if l.Level <= ERROR {
		l.logger.Println("[ERROR]", fmt.Sprintf("%s %v", msg, formatKeysAndValues(keysAndValues...)))
	}
}
func (l *Logger) Debug(msg string, keysAndValues ...interface{}){
	if l.Level <= DEBUG {
		l.logger.Println("[DEBUG]", fmt.Sprintf("%s %v", msg, formatKeysAndValues(keysAndValues...)))
	}
}


func formatKeysAndValues(keysAndValues ...interface{}) string {
	data := make(map[string]interface{})
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			data[fmt.Sprint(keysAndValues[i])] = keysAndValues[i+1]
		} else {
			data[fmt.Sprint(keysAndValues[i])] = "MISSING"
		}
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("Error formatting data: %v", err)
	}
	return string(jsonData)
}