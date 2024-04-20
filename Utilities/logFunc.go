package Utilities

import (
	"VideoWeb/logrusLog"
)

func WriteErrLog(funcName, message string) {
	logrusLog.Log.WithFields(map[string]interface{}{
		"type":     "Error",
		"function": funcName,
	}).Errorf(message)
}

func WriteInfoLog(funcName, message string) {
	logrusLog.Log.WithFields(map[string]interface{}{
		"type":     "Info",
		"function": funcName,
	}).Infof(message)
}
