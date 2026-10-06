package main

import (
	"Service/Agent/AgentService"
	"Service/Agent/Clients"
	"Service/Tools/sqlTools"
	"Service/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("出错: .env 文件未被找到")
	}

	//---初始化数据库---
	db := sqlTools.InitSQLConnect()

	//---配置服务器---
	serve := gin.Default()
	serve.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	server.StartServer(serve, db)             //面向前端服务
	AgentService.AgentSQLService(serve, db)   //面向Agent服务
	Clients.HandleConnectClientService(serve) //中转站网关Tool

	//---启动服务器---
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	baseUrl := fmt.Sprintf("%s%s", os.Getenv("SERVER_URL"), addr)
	fmt.Printf("服务端启动在 %s\n", baseUrl)
	if err := http.ListenAndServe(addr, serve); err != nil {
		log.Fatal("服务器启动失败:", err)
	}

}
