package Utilities

import (
	"VideoWeb/logrusLog"
)

func WriteLog(funcName, message string) {
	logrusLog.Log.WithField("function", funcName).Errorf(message)
	//fmt.Println(logrusLog.Log.WithField("funcName", funcName).Data)
}
