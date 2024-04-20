package test

import (
	"VideoWeb/Utilities"
	"VideoWeb/config"
	"VideoWeb/logrusLog"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Println("Log is:", logrusLog.Log)
	err := config.InitConfig("D:\\Go\\WorkSpace\\src\\Go_Project\\VideoWebFromUbuntu-22.04LTS\\VideoWeb\\config\\config.yaml")
	if err != nil {
		t.Log("err:", err)
		return
	}
	Utilities.WriteLog("TestLog", "This is a test log")

}
