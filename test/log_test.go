package test

import (
	"VideoWeb/Utilities/logf"
	"VideoWeb/config"
	"testing"
)

func TestLog(t *testing.T) {
	config.InitConfig("D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWeb\\VideoWeb\\config\\config.yaml")

	logf.WriteErrLog("TestLog", "This is a test log")
	logf.WriteInfoLog("TestLog", "This is a test log")

}
