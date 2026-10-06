package image

import (
	"Service/Tools/JWT"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"log"
	"net/http"
	"time"
)

func HandleImageQiNiu(r *gin.Engine, db *sql.DB) {
	t := r.Group("/api/image/upload")
	t.GET("/token", JWT.AuthMiddleware(), func(c *gin.Context) {
		_, _, bucket, domain, mac := getConfig()

		log.Printf("七牛配置: Bucket=%q, Domain=%q", bucket, domain)

		putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(2*time.Hour))
		if err != nil {
			log.Println("NewPutPolicy 失败:", err)
			c.JSON(500, gin.H{"error": "Bucket配置错误: " + err.Error()})
			return
		}
		upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
		if err != nil {
			log.Println("GetUpToken 失败:", err)
			c.JSON(500, gin.H{"error": "签名失败: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"uptoken": upToken,
			"domain":  domain,
		})
	})
	t.POST("/result", JWT.AuthMiddleware(), func(c *gin.Context) {
		var req Response
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("收到前端回传: key=%s hash=%s", req.Key, req.Hash)
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}
