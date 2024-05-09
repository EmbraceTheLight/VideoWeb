package Utilities

import (
	"VideoWeb/define"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetIPInfo 获得用户的IP信息
func GetIPInfo(IP string) (*define.IPInfo, error) {
	searchURL := "http://ip-api.com/json/" + IP + "?lang=zh-CN"
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipInfo define.IPInfo
	err = json.Unmarshal(out, &ipInfo)
	if err != nil {
		fmt.Println("err in Unmarshal:", err)
		return nil, err
	}

	return &ipInfo, nil

}

// GetMyPublicIP 通过访问http://myexternalip.com/raw获取公网ip
func GetMyPublicIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	return string(content)
}
