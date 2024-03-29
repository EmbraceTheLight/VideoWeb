package main

import (
	"VideoWeb/DAO"
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
	r := gin.Default()
	err := DAO.InitMySQL()
	if err != nil {
		fmt.Println("err in InitMySQL:", err)
	}

	//解决跨域问题，注册全局中间件
	r.Use(cors.Default())
	routers.CollectRouter(r)
	log.Fatal(r.Run("0.0.0.0:8080"))
}
