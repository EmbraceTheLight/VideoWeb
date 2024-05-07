package test

import (
	"VideoWeb/Utilities/logf"
	"VideoWeb/config"
	"testing"
)

func TestLog(t *testing.T) {
	err := config.InitConfig("D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWebFromUbuntu-22.04LTS\\VideoWeb\\config\\config.yaml")
	if err != nil {
		t.Log("err:", err)
		return
	}
	logf.WriteErrLog("TestLog", "This is a test log")
	logf.WriteInfoLog("TestLog", "This is a test log")

}
