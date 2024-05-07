package logf

import (
	"VideoWeb/logrusLog"
)

func WriteErrLog(funcName, message string) {
	logrusLog.Log.WithFields(map[string]interface{}{
		"function": funcName,
	}).Errorf(message)
}

func WriteInfoLog(funcName, message string) {
	logrusLog.Log.WithFields(map[string]interface{}{
		"function": funcName,
	}).Infof(message)
}
