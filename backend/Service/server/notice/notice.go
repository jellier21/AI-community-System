package notice

import (
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleNotice(r *gin.Engine, db *sql.DB) {
	g := r.Group("/api/notice")

	// 发布公告（只有 leader 能发）
	g.POST("/publish", JWT.AuthMiddleware(), func(c *gin.Context) {
		leaderAccount := c.GetString("userID")
		communityUid := c.PostForm("community_uid")
		title := c.PostForm("title")
		content := c.PostForm("content")

		if communityUid == "" || title == "" || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失"})
			return
		}

		// 校验是 leader
		err := sqlTools.CheckLeaderOfCommunity(leaderAccount, communityUid, db)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "只有负责人能发布公告"})
			return
		}

		err = sqlTools.PublishNotice(communityUid, title, content, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "发布成功"})
	})

	// 查公告列表
	g.GET("/list", JWT.AuthMiddleware(), func(c *gin.Context) {
		communityUid := c.Query("community_uid")
		list, err := sqlTools.GetNoticeList(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"info": list})
	})

	// 查最新一条（home 页用）
	g.GET("/latest", JWT.AuthMiddleware(), func(c *gin.Context) {
		communityUid := c.Query("community_uid")
		n, err := sqlTools.GetLatestNotice(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"info": n})
	})
}
