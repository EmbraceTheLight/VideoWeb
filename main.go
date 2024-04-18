package main

import (
	"VideoWeb/DAO"
	"VideoWeb/Utilities"
	"VideoWeb/routers"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// @title           VideoWeb
// @version         1.0
// @description     This is a VideoWeb API
func main() {
	//time.Sleep(5 * time.Second)
	r := gin.Default()
	err := DAO.InitDB()
	defer DAO.RDB.Close()
	if err != nil {
		fmt.Println("err in InitDB:", err)
	}

	go Utilities.HardDelete()
	//解决跨域问题，注册全局中间件
	r.Use(cors.Default())
	routers.CollectRouter(r)
	log.Fatal(r.Run(":51233"))
}
