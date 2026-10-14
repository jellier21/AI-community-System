package share

import (
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandleShare(r *gin.Engine, db *sql.DB) {
	g := r.Group("/api/share")

	// 发布
	g.POST("/publish", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		communityUid := c.PostForm("community_uid")
		content := c.PostForm("content")
		images := c.PostForm("images") // 暂时为空

		if communityUid == "" || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
			return
		}
		err := sqlTools.PublishShare(account, communityUid, content, images, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "发布成功"})
	})

	// 我自己的动态
	g.GET("/my", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		list, err := sqlTools.GetMyShares(account, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"info": list})
	})

	// 本社区所有动态（管理员）
	g.GET("/list", JWT.AuthMiddleware(), func(c *gin.Context) {
		communityUid := c.Query("community_uid")
		list, err := sqlTools.GetCommunityShares(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"info": list})
	})

	// 全站最新动态（首页"为您推荐"用）
	g.GET("/latest", JWT.AuthMiddleware(), func(c *gin.Context) {
		list, err := sqlTools.GetLatestShares(3, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"info": list})
	})

	// 删除
	g.POST("/delete", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		idStr := c.PostForm("id")
		id, _ := strconv.Atoi(idStr)
		err := sqlTools.DeleteShare(id, account, db)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权删除"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	g.POST("/update", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		id, _ := strconv.Atoi(c.PostForm("id"))
		content := c.PostForm("content")
		images := c.PostForm("images")
		err := sqlTools.UpdateShare(id, account, content, images, db)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
	})
}
