package logger

import (
	"fmt"
	"log"
)

type Logger struct{}

func (l Logger) Info(msg string, keysAndValues ...interface{}){
	log.Printf("[INFO] %s %v\n", msg, formatKeysAndValues(keysAndValues...))
}
func(l Logger) Error(msg string, keysAndValues ...interface{}){
	log.Printf("[ERROR] %s %v\n", msg, formatKeysAndValues(keysAndValues...))
}

func formatKeysAndValues(keysAndValues ...interface{}) string {
	var formatted []interface{}
	for i := 0; i < len(keysAndValues); i += 2{
		if i+1 < len(keysAndValues) {
			formatted = append(formatted, fmt.Sprintf("%v=%v", keysAndValues[i], keysAndValues[i+1]))		
		} else {
			formatted = append(formatted, fmt.Sprintf("%v=MISSING", keysAndValues[i]))
		}
	}
	return fmt.Sprint(formatted...)
}