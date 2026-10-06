package Clients

import (
	"Service/Agent/Clients/mapreading"
	"Service/Agent/Clients/medical"
	"Service/Agent/Clients/teacher"
	"Service/Tools/JWT"
	"github.com/gin-gonic/gin"
)

func HandleConnectClientService(r *gin.Engine) {
	agent := r.Group("/api/agent/")
	agent.GET("/test", func(c *gin.Context) {
		res, err := medical.TestClient(c.DefaultQuery("message", ""))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"result": res})
	})

	agent.GET("/text", func(c *gin.Context) {
		res, err := medical.TextAgentClient(c.DefaultQuery("message", ""))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"result": res})
	})
	agent.Group("/medical").POST("/diagnosis", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.PostForm("account")
		question := c.PostForm("question")
		imageFile, fileErr := c.FormFile("file")
		if fileErr != nil && question == "" {
			c.JSON(400, gin.H{"error": "需要文字或图片"})
			return
		}
		response, err := medical.MedicalAgentClient(account, question, imageFile) // 文件可能为 nil
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"response": response})
	})
	agent.Group("/teacher").POST("/ask", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.PostForm("account")
		question := c.PostForm("question")
		imageFile, fileErr := c.FormFile("file")
		if account == "" {
			c.JSON(400, gin.H{"error": "账号为必填项"})
			return
		}
		if fileErr != nil && question == "" {
			c.JSON(400, gin.H{"error": "需要文字或图片"})
			return
		}
		response, err := teacher.TeacherAgentClient(account, question, imageFile)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"response": response})
	})
	agent.Group("/map").POST("/recognize", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.PostForm("account")
		imageFile, fileErr := c.FormFile("file")
		if fileErr != nil {
			c.JSON(400, gin.H{"error": "请上传图片"})
			return
		}
		response, err := mapreading.MapReadingClient(account, imageFile)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"response": response})
	})

}
