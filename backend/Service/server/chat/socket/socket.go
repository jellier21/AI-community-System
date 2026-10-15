package socket

import (
	"Service/Tools/JWT"
	"Service/Tools/sqlTools"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

func NewSocket(db *sql.DB) *Hub {
	//---Hub---
	hub := &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *ClientRegistration),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *RoomMessage, 256),
		DB:         db,
	}
	return hub
}

func HandleSocket(r *gin.Engine, db *sql.DB) {
	hub := NewSocket(db)
	go hub.Run(db)
	r.GET("/ws", JWT.AuthMiddleware(), func(c *gin.Context) {
		_account_user_ := c.GetString("userID")
		if _account_user_ == "" {
			c.JSON(400, gin.H{"error": "无法识别当前用户"})
			return
		}

		// 限制账号长度（防止过长账号破坏 UI）
		if len([]rune(_account_user_)) > 20 {
			_account_user_ = string([]rune(_account_user_)[:20])
		}

		roomUid := c.Request.URL.Query().Get("roomUid")
		if roomUid == "" {
			c.JSON(400, gin.H{"error": "RoomUid is required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			c.JSON(500, gin.H{"error": "升级连接失败"})
			log.Printf("升级失败: %v", err)
			return
		}

		//Agent := NewClient(conn, _account_user_)
		client := &Client{
			Account:     _account_user_,
			Conn:        conn,
			CurrentRoom: roomUid,
			Send:        make(chan []byte, sendBufferSize),
			Hub:         hub,
		}

		//加入默认房间
		hub.Register <- &ClientRegistration{
			Client:  client,
			RoomUid: roomUid,
		}
		go client.writePump() // 启动写协程

		// 心跳设置
		conn.SetReadDeadline(time.Now().Add(ReadTimeout))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(ReadTimeout))
			return nil
		})
		defer func() {
			hub.Unregister <- client
			close(client.Send)
			conn.Close()
		}()
		// 读循环
		for {
			_, rawMessage, err := conn.ReadMessage()
			if err != nil {
				break
			}
			conn.SetReadDeadline(time.Now().Add(ReadTimeout))
			msg, err := ParseMessage(rawMessage)
			if err != nil {
				msg = NewMessage(_account_user_, string(rawMessage))
			}
			msg.Account = _account_user_
			// 查昵称和头像
			msg.Nickname = sqlTools.GetMemberDisplayName(client.CurrentRoom, _account_user_, client.Hub.DB)
			msg.Avatar = sqlTools.GetUserAvatarByAccount(_account_user_, client.Hub.DB)
			data, _ := msg.ToJSON()

			// 房间内广播
			hub.Broadcast <- &RoomMessage{
				RoomUid: client.CurrentRoom,
				Data:    data,
				Sender:  client,
			}
		}
	})
	w := r.Group("/api/ws")
	// 获取历史消息API
	w.GET("/messages/history", JWT.AuthMiddleware(), func(c *gin.Context) {
		roomUid := c.Query("room_uid")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

		messages, err := sqlTools.GetRoomMessages(roomUid, page, pageSize, db)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "获取历史消息失败",
				"data":    []interface{}{}, // 返回空数组而不是null
			})
			return
		}

		if messages == nil {
			messages = []sqlTools.MessageDB{} // 确保不返回null
		}

		c.JSON(200, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    messages,
		})
	})

	// 获取最近消息API（用户进入房间时）
	w.GET("/messages/recent", JWT.AuthMiddleware(), func(c *gin.Context) {
		roomUid := c.Query("room_uid")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		messages, err := sqlTools.GetRecentMessages(roomUid, limit, db)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "获取最近消息失败",
				"data":    []interface{}{}, // 返回空数组而不是null
			})
			return
		}

		if messages == nil {
			messages = []sqlTools.MessageDB{} // 确保不返回null
		}

		c.JSON(200, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    messages,
		})
	})

	w.GET("/api/onlinetool/count", func(c *gin.Context) {
		mu.Lock()
		count := len(clients)
		mu.Unlock()

		c.Writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(c.Writer, `{"online":%d}`, count)
	})
}
