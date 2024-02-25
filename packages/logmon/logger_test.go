package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"errors"
	"strings"
	"testing"
)

func setupLoggerWithBuffer(level LogLevel, prefix string) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, level, "")
	return logger, &buf
}

func TestInfoLevelLogging(t *testing.T) {
	logger, buf := setupLoggerWithBuffer(INFO, "")
	logger.Info("Info message", "key", "value")
	expectedSubstring := `"key":"value"`
	if !strings.Contains(buf.String(), expectedSubstring) {
		t.Errorf("Expected INFO log to contain %s, got %s", expectedSubstring, buf.String())
	}
}

func TestErroLevelLogging(t *testing.T) {
	logger, buf := setupLoggerWithBuffer(ERROR, "")
	logger.Info("Message should not be logged", "info key", "infoValue")
	logger.Error("Error occured", "errorKey", "errorValue")
	expectedSubString := `"errorKey":"errorValue"`
	if !strings.Contains(buf.String(), expectedSubString) {
		t.Errorf("Expected ERROR log to contain %s, got %s", expectedSubString, buf.String())
	}
	unexpectedSubstring := `"infoKey":"infoValue"`
	if strings.Contains(buf.String(), unexpectedSubstring) {
		t.Errorf("Expected ERROR log to be present, but got INFO log instead")
	}
}

func TestDebugLevelLogging(t *testing.T) {
	logger, buf := setupLoggerWithBuffer(DEBUG, "")
	logger.Debug("Debugging", "debugKey", "debugValue")
	expectedSubstring := `"debugKey":"debugValue"`
	if !strings.Contains(buf.String(), expectedSubstring) {
		t.Errorf("Expected DEBUG log to contain %s, got %s", expectedSubstring, buf.String())
	}
}

func isolatedJSON(logOutput string) (map[string]interface{}, error) {
	startIndex := strings.Index(logOutput, "{")
	if startIndex == -1 {
		return nil, errors.New("JOSN object not found in log output")
	}
	jsonPart := logOutput[startIndex:]
	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonPart), &logData); err != nil {
		return nil, fmt.Errorf("failed to parse log as JSON: %w", err)
	}
	return logData, nil
}

func VerifyLogFormat(t *testing.T, logOutput, key, expectedValue string) {
	logData, err := isolatedJSON(logOutput)
	if err != nil {
		t.Fatalf(err.Error())
	}

	val, ok := logData[key]
	if !ok {
		t.Errorf("Log does not contain the expected key: %s", key)
	} else {
		valStr := fmt.Sprintf("%v", val)
		if valStr != expectedValue {
			t.Errorf("Expected value for key is %q is %q, got %q", key, expectedValue, valStr)
		}
	}
}

func TestINFOStructured(t *testing.T) {
	logger, buf := setupLoggerWithBuffer(INFO, "")
	logger.Info("INFO Log test", "userID", "098765", "action", "login")
	logOutput := buf.String()
	VerifyLogFormat(t, logOutput, "userID", "098765")
	VerifyLogFormat(t, logOutput, "action", "login")
}

func TestERRORStructured(t *testing.T) {
	logger, buf := setupLoggerWithBuffer(ERROR, "")
	logger.Error("Error Log test", "userID", "098765", "action", "error")
	logOutput := buf.String()
	VerifyLogFormat(t, logOutput, "userID", "098765")
	VerifyLogFormat(t, logOutput, "action", "error")
}

func TestDEBUGStructured(t *testing.T) {
	logger, buf := setupLoggerWithBuffer(DEBUG, "")
	logger.Debug("Debug Log test", "userID", "12345", "action", "debug")
	logOutput := buf.String()
	VerifyLogFormat(t, logOutput, "userID", "12345")
	VerifyLogFormat(t, logOutput, "action", "debug")
}
