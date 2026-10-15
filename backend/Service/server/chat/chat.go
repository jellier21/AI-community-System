package chat

import (
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"Service/server/chat/socket"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func WebsocketServer(r *gin.Engine, db *sql.DB) {
	socket.HandleSocket(r, db)
	group := r.Group("/api/chat/group")
	group.GET("/list", JWT.AuthMiddleware(), func(c *gin.Context) {
		_account_ := c.GetString("userID")
		var RoomList []sqlTools.RoomDB
		RoomList, err := sqlTools.GetChatroomLists(_account_, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"code":      200,
			"room_list": RoomList,
		})
	})

	// 发起私聊：输入对方账号
	group.POST("/private/create", JWT.AuthMiddleware(), func(c *gin.Context) {
		myAccount := c.GetString("userID")
		friendAccount := strings.TrimSpace(c.PostForm("friend_account"))
		if myAccount == "" || friendAccount == "" {
			c.JSON(400, gin.H{"code": 400, "error": "参数缺失"})
			return
		}
		if myAccount == friendAccount {
			c.JSON(400, gin.H{"code": 400, "error": "不能和自己私聊"})
			return
		}
		exists, nick := sqlTools.IsUserExists(friendAccount, db)
		if !exists {
			c.JSON(400, gin.H{"code": 400, "error": "对方账号不存在"})
			return
		}
		roomUid := sqlTools.PrivateRoomUid(myAccount, friendAccount)
		var existing string
		err := db.QueryRow("SELECT uid FROM room WHERE uid = ?", roomUid).Scan(&existing)
		log.Printf("私聊 roomUid=%s, existing err=%v", roomUid, err)
		if err == sql.ErrNoRows {
			err = sqlTools.CreatePrivateRoom(myAccount, friendAccount, nick, db)
			if err != nil {
				log.Printf("创建私聊房间失败: %v", err)
				c.JSON(500, gin.H{"code": 500, "error": err.Error()})
				return
			}
		}
		c.JSON(200, gin.H{"code": 200, "room_uid": roomUid, "friend_nick": nick})
	})

	// 创建群聊
	group.POST("/create", JWT.AuthMiddleware(), func(c *gin.Context) {
		myAccount := c.GetString("userID")
		groupName := strings.TrimSpace(c.PostForm("group_name"))
		membersStr := c.PostForm("members")
		if myAccount == "" || groupName == "" {
			c.JSON(400, gin.H{"code": 400, "error": "群名不能为空"})
			return
		}
		members := strings.Split(membersStr, ",")
		roomUid, err := sqlTools.CreateGroupRoom(myAccount, groupName, members, db)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200, "room_uid": roomUid})
	})

	// 好友列表
	group.GET("/friends", JWT.AuthMiddleware(), func(c *gin.Context) {
		myAccount := c.GetString("userID")
		list, err := sqlTools.GetFriendList(myAccount, db)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200, "friends": list})
	})

	// 删除好友
	group.POST("/friend/delete", JWT.AuthMiddleware(), func(c *gin.Context) {
		myAccount := c.GetString("userID")
		roomUid := c.PostForm("room_uid")
		if roomUid == "" {
			c.JSON(400, gin.H{"code": 400, "error": "房间号不能为空"})
			return
		}
		// 校验当前用户确实在这个房间里，防止删别人的私聊
		var cnt int
		err := db.QueryRow("SELECT COUNT(*) FROM `groups` WHERE uid = ? AND account = ?", roomUid, myAccount).Scan(&cnt)
		if err != nil || cnt == 0 {
			c.JSON(403, gin.H{"code": 403, "error": "无权删除该会话"})
			return
		}
		err = sqlTools.DeletePrivateRoom(roomUid, db)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200})
	})
}
