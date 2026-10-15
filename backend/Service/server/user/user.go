package user

import (
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func StartUserServeAPI(r *gin.Engine, db *sql.DB) {
	user := r.Group("/api/user")
	user.POST("/register", func(c *gin.Context) {
		// 尝试解析multipart表单
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			// 如果multipart失败，尝试解析普通表单
			if err := c.Request.ParseForm(); err != nil {
				log.Printf("解析表单失败: %v", err)
				c.JSON(400, gin.H{
					"code":  400,
					"error": "Failed to parse form data: " + err.Error(),
				})
				return
			}
		}

		account := c.PostForm("account")
		password := c.PostForm("password")
		if account == "" || password == "" {
			c.JSON(400, gin.H{
				"code":  400,
				"error": "Account and password are required",
			})
			log.Printf("Received empty account or password")
			return
		}
		log.Printf("收到User: %s的注册申请，密码： %s", account, password)

		// 检查数据是否为空

		err := sqlTools.RegisterTools(account, password, db)
		if err != nil {
			c.JSON(500, gin.H{
				"code":  500,
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Registration successful",
		})
	})
	user.POST("/login", func(c *gin.Context) {
		account := c.PostForm("account")
		password := c.PostForm("password")
		if account == "" || password == "" {
			c.JSON(400, gin.H{
				"code":  400,
				"error": "Account and password are required",
			})
			log.Printf("Received empty account or password")
			return
		}
		log.Printf("收到User: %s的登录申请，密码： %s", account, password)

		token, err := sqlTools.LoginTools(account, password, db)
		if err != nil {
			c.JSON(500, gin.H{
				"code":  500,
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Login successful",
			"token":   token,
		})
	})
	query := user.Group("/query")
	{
		query.GET("/nick", JWT.AuthMiddleware(), func(c *gin.Context) {
			account := c.Query("account")
			nickname, err := sqlTools.QueryUserNickname(account, db)
			if err != nil {
				c.JSON(500, gin.H{
					"code":  500,
					"error": "Failed to find nickname",
				})
				return
			}
			c.JSON(200, gin.H{
				"code":     200,
				"message":  "Nickname found successfully",
				"nickname": nickname,
			})
		})
		query.GET("/avatar", JWT.AuthMiddleware(), func(c *gin.Context) {
			account := c.Query("account")
			avatar, err := sqlTools.QueryUserAvatar(account, db)
			if err != nil {
				c.JSON(500, gin.H{"code": 500, "error": "查询头像失败"})
				return
			}
			c.JSON(200, gin.H{"code": 200, "avatar": avatar})
		})
		query.GET("/profile", JWT.AuthMiddleware(), func(c *gin.Context) {
			account := c.Query("account")
			profile, err := sqlTools.QueryUserProfile(account, db)
			if err != nil {
				c.JSON(500, gin.H{"code": 500, "error": "查询用户资料失败"})
				return
			}
			c.JSON(200, gin.H{"code": 200, "profile": profile})
		})

	}
	amend := user.Group("/amend")
	amend.POST("/nick", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		nickname := c.PostForm("nickname")
		if account == "" {
			c.JSON(400, gin.H{
				"code":  400,
				"error": "Account is required",
			})
			log.Printf("Received empty account")
			return
		}
		if nickname == "" {
			c.JSON(400, gin.H{
				"code":  400,
				"error": "Nickname is required",
			})
			log.Printf("Received empty nickname")
			return
		}
		log.Printf("收到User: %s的昵称修改申请，昵称： %s", account, nickname)

		err := sqlTools.AmendNickByAccount(nickname, account, db)
		if err != nil {
			log.Printf("Nickname amendment failed: %v", err)
			c.JSON(500, gin.H{
				"code":  500,
				"error": "Nickname amendment failed",
			})
			return
		}
		log.Printf("Nickname amendment successful: %s", account)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Nickname amendment successful",
		})
	})
	amend.POST("/password", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		password := c.PostForm("password_new")
		password2 := c.PostForm("password_old")
		if account == "" || password == "" {
			c.JSON(400, gin.H{
				"code":  400,
				"error": "Account and password are required",
			})
			log.Printf("Received empty account or password")
			return
		}
		log.Printf("收到User: %s的密码修改申请", account)

		err := sqlTools.AmendPasswordByAccount(account, password, password2, db)
		if err != nil {
			c.JSON(500, gin.H{
				"code":  500,
				"error": "Password amendment failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Password amendment successful",
		})
	})
	amend.POST("/signature", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		signature := c.PostForm("signature")
		if account == "" || signature == "" {
			c.JSON(400, gin.H{
				"code":  400,
				"error": "Account and signature are required",
			})
			log.Printf("Received empty account or signature")
			return
		}
		log.Printf("收到User: %s的签名修改申请", account)

		err := sqlTools.AmendMySignatureByAccount(signature, account, db)
		if err != nil {
			c.JSON(500, gin.H{
				"code":  500,
				"error": "Signature amendment failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Signature amendment successful",
		})
	})
	amend.POST("/avatar", JWT.AuthMiddleware(), func(c *gin.Context) {
		account := c.GetString("userID")
		avatar := c.PostForm("avatar")
		if account == "" {
			c.JSON(400, gin.H{"code": 400, "error": "账号不能为空"})
			return
		}
		err := sqlTools.AmendAvatarByAccount(avatar, account, db)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "error": "更新头像失败"})
			return
		}
		c.JSON(200, gin.H{"code": 200, "message": "头像更新成功"})
	})

}
