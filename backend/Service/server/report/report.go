package report

import (
	reportClient "Service/Agent/Clients/report"
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleReport(r *gin.Engine, db *sql.DB) {
	g := r.Group("/api/report")

	// 生成 AI 周报：取本社区最近一周公告 + 医疗咨询，发给 AI 总结
	g.POST("/generate", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")

		// 1. 拿社区号和社区名
		communityUid, err := sqlTools.GetCommunityUID(account, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未查询到所属社区"})
			return
		}
		communityInfo, err := sqlTools.GetCommunityInfo(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询社区信息失败"})
			return
		}

		// 2. 查最近 7 天公告 + 医疗咨询
		notices, err := sqlTools.GetWeeklyNotices(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询公告失败"})
			return
		}
		medicals, err := sqlTools.GetWeeklyMedical(communityUid, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询医疗记录失败"})
			return
		}

		// 3. 转发 AI 端生成周报
		res, err := reportClient.WeeklyReportClient(communityInfo.Name, notices, medicals)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"report": res})
	})
}
