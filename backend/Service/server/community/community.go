package community

import (
	"Service/Tools/JWT"
	"Service/Tools/publish"
	"Service/Tools/sqlTools"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func HandleCommunity(r *gin.Engine, db *sql.DB) {
	community := r.Group("/api/community")
	community.POST("/create", JWT.AuthMiddleware(), func(c *gin.Context) {
		_leader_ := c.PostForm("account")
		_name_ := c.PostForm("name")
		_description_ := c.PostForm("description")
		//log.Printf("Creating community: %s, %s, %s", _leader_, _name_, _description_)

		Uid := publish.MakeID(_leader_, time.Now().UnixMilli())

		err := sqlTools.BuildCommunitySQL(_leader_, _name_, _description_, Uid, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Community created successfully"})
	})
	community.POST("/delete", JWT.AuthMiddleware(), func(c *gin.Context) {
		_account_ := c.PostForm("account")
		err := sqlTools.QuitCommunity(_account_, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Quit community successfully"})
	})
	query := community.Group("/query")
	query.GET("/have", JWT.AuthMiddleware(), func(c *gin.Context) {
		_account_ := c.Query("account")
		hc, err := sqlTools.UserWheatherHaveCommunity(_account_, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"info": hc,
		})
	})
	query.GET("/uid", JWT.AuthMiddleware(), func(c *gin.Context) {
		_account_ := c.Query("account")
		uid, err := sqlTools.GetCommunityUID(_account_, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"info": uid,
		})
	})
	query.GET("/isleader", JWT.AuthMiddleware(), func(c *gin.Context) {
		_account_ := c.Query("account")
		err := sqlTools.UserWheatherIsLeader(_account_, db)
		if err != nil {
			c.JSON(200, gin.H{
				"info": "user",
			})
			return
		}
		c.JSON(200, gin.H{
			"info": "leader",
		})
	})
	member := query.Group("/member")
	member.GET("/list", JWT.AuthMiddleware(), func(c *gin.Context) {
		_community_ := c.Query("uid")
		memberList, err := sqlTools.GetMyMemberList(_community_, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"info": memberList,
		})
	})
	member.POST("/remove", JWT.AuthMiddleware(), func(c *gin.Context) {
		// 当前登录人账号（JWT 中间件写入）
		leaderAccount := c.GetString("userID")
		communityUid := c.PostForm("community_uid")
		accountsStr := c.PostForm("accounts")
		log.Printf("Removing members: %s from community %s by %s", accountsStr, communityUid, leaderAccount)

		if communityUid == "" || accountsStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失"})
			return
		}

		// 校验当前登录人确实是这个社区的 leader
		err := sqlTools.CheckLeaderOfCommunity(leaderAccount, communityUid, db)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权操作该社区"})
			return
		}

		// 逗号分隔 → 切片，并剔除空串和 leader 自己
		accountList := strings.Split(accountsStr, ",")
		filtered := make([]string, 0, len(accountList))
		for _, a := range accountList {
			a = strings.TrimSpace(a)
			if a == "" || a == leaderAccount {
				continue
			}
			filtered = append(filtered, a)
		}
		if len(filtered) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "无可删除成员"})
			return
		}

		err = sqlTools.RemoveMembers(communityUid, filtered, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	// leader 添加成员
	member.POST("/add", JWT.AuthMiddleware(), func(c *gin.Context) {
		leaderAccount := c.GetString("userID")
		communityUid := c.PostForm("community_uid")
		targetAccount := c.PostForm("account")
		memberName := c.PostForm("member_name")

		if communityUid == "" || targetAccount == "" || memberName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失"})
			return
		}

		// 校验是 leader
		err := sqlTools.CheckLeaderOfCommunity(leaderAccount, communityUid, db)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权操作该社区"})
			return
		}

		err = sqlTools.AddMember(communityUid, targetAccount, memberName, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
	})

	// 成员通过社区号加入
	member.POST("/join", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		communityUid := c.PostForm("community_uid")

		if communityUid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入社区号"})
			return
		}

		err := sqlTools.JoinCommunityByUid(communityUid, account, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "加入成功"})
	})

}
