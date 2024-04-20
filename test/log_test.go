package test

import (
	"VideoWeb/Utilities"
	"VideoWeb/config"
	"testing"
)

func TestLog(t *testing.T) {
	err := config.InitConfig("D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWebFromUbuntu-22.04LTS\\VideoWeb\\config\\config.yaml")
	if err != nil {
		t.Log("err:", err)
		return
	}
	Utilities.WriteErrLog("TestLog", "This is a test log")
	Utilities.WriteInfoLog("TestLog", "This is a test log")

}
