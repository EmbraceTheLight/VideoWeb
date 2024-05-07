package main

import (
	"VideoWeb/DAO"
	"VideoWeb/Utilities"
	"VideoWeb/config"
	"VideoWeb/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// @title           VideoWeb
// @version         1.0
// @description     This is a VideoWeb API
func main() {
	config.InitConfig("")

	DAO.InitDBs()
	defer DAO.RDB.Close()

	//执行后台定时任务：删除软删除记录
	go Utilities.HardDelete()

	r := gin.Default()
	//解决跨域问题，注册全局中间件
	r.Use(cors.Default())
	routers.CollectRouter(r)
	log.Fatal(r.Run(":51233"))
}
