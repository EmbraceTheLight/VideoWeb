package test

import (
	"VideoWeb/Utilities"
	"fmt"
	"testing"
)

func TestGetIPInfo(t *testing.T) {
	IP := Utilities.GetMyPublicIP()
	info, _ := Utilities.GetIPInfo(IP)
	info2, _ := Utilities.GetIPInfo("47.242.47.233")

	fmt.Println(info)
	fmt.Println(info2)
}
