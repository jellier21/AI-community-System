package advance

import (
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleAdvance(r *gin.Engine, db *sql.DB) {
	g := r.Group("/api/advance")

	// 用户提交建议
	g.POST("/submit", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		communityUid := c.PostForm("community_uid")
		title := c.PostForm("title")
		content := c.PostForm("content")

		if communityUid == "" || title == "" || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失"})
			return
		}

		err := sqlTools.SubmitAdvance(account, communityUid, title, content, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "提交成功"})
	})

	// leader 查看建议列表
	g.GET("/list", JWT.AuthMiddleware(), func(c *gin.Context) {
		communityUid := c.Query("community_uid")
		list, err := sqlTools.GetAdvanceList(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"info": list})
	})
}
