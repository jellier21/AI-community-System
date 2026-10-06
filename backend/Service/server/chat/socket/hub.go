package socket

import (
	"Service/Tools/sqlTools"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"log"
)

// Run 启动 Hub 的主事件循环
func (h *Hub) Run(db *sql.DB) {
	for {
		select {
		case reg := <-h.Register:
			h.handleRegister(reg, db)
		case client := <-h.Unregister:
			h.handleUnregister(client, db)
		case msg := <-h.Broadcast:
			h.handleBroadcast(msg)
		}
	}
}

// handleRegister 处理客户端加入房间
func (h *Hub) handleRegister(reg *ClientRegistration, db *sql.DB) {
	userNick, _ := sqlTools.QueryUserNickname(reg.Client.Account, db)
	// 1. 获取或创建房间
	h.RoomsMu.Lock()
	room, exists := h.Rooms[reg.RoomUid]
	if !exists {
		room = &Room{
			RoomUid: reg.RoomUid,
			Clients: make(map[*Client]bool),
		}
		h.Rooms[reg.RoomUid] = room
		log.Printf("创建房间: %s", reg.RoomUid)
	}
	h.RoomsMu.Unlock()
	// 2. 将客户端加入房间
	room.RoomMu.Lock()
	room.Clients[reg.Client] = true
	clientCount := len(room.Clients)
	room.RoomMu.Unlock()
	// 3. 设置客户端的当前房间
	reg.Client.CurrentRoom = reg.RoomUid
	log.Printf("[%s] 加入房间: %s (房间人数: %d)",
		userNick, reg.RoomUid, clientCount)
	// 4. 广播加入通知
	joinMsg := SystemMessage("join",
		fmt.Sprintf("%s 加入了房间 %s (共 %d 人)",
			userNick, reg.RoomUid, clientCount))
	if data, err := joinMsg.ToJSON(); err == nil {
		h.broadcastToRoom(reg.RoomUid, data, reg.Client)
	}
}

// handleUnregister 处理客户端离开
func (h *Hub) handleUnregister(client *Client, db *sql.DB) {
	userNick, _ := sqlTools.QueryUserNickname(client.Account, db)
	roomName := client.CurrentRoom
	if roomName == "" {
		return
	}
	h.RoomsMu.RLock()
	room, exists := h.Rooms[roomName]
	h.RoomsMu.RUnlock()
	if !exists {
		return
	}
	// 从房间移除客户端
	room.RoomMu.Lock()
	delete(room.Clients, client)
	remaining := len(room.Clients)
	room.RoomMu.Unlock()
	log.Printf("[%s] 离开房间 %s (剩余: %d)",
		userNick, roomName, remaining)
	// 广播离开通知
	leaveMsg := SystemMessage("leave",
		fmt.Sprintf("%s 离开了房间 %s (剩余 %d 人)",
			userNick, roomName, remaining))
	if data, err := leaveMsg.ToJSON(); err == nil {
		h.broadcastToRoom(roomName, data, nil)
	}
	// 空房间自动清理
	if remaining == 0 {
		h.RoomsMu.Lock()
		delete(h.Rooms, roomName)
		h.RoomsMu.Unlock()
		log.Printf("房间 %s 已无人，自动清理", roomName)
	}
}

// handleBroadcast 处理房间广播
func (h *Hub) handleBroadcast(msg *RoomMessage) {
	h.RoomsMu.RLock()
	room, exists := h.Rooms[msg.RoomUid]
	h.RoomsMu.RUnlock()
	if !exists {
		return
	}
	room.RoomMu.Lock()
	defer room.RoomMu.Unlock()

	// 解析消息内容用于存储
	var message Message
	if err := json.Unmarshal(msg.Data, &message); err == nil {
		// 异步保存消息到数据库（避免阻塞广播）
		go func(db *sql.DB) {
			if err := sqlTools.SaveMessage(
				msg.RoomUid,
				message.Account,
				message.Content,
				db,
			); err != nil {
				log.Println("消息持久化失败:", err)
			}
		}(h.DB) // 需要在Hub结构中添加DB字段
	}

	// 广播给其他客户端
	for client := range room.Clients {
		if client == msg.Sender {
			continue
		}
		client.SendMessage(msg.Data)
	}
}

// broadcastToRoom 向房间内广播（不经过 channel，直接调用）
func (h *Hub) broadcastToRoom(roomName string, data []byte, sender *Client) {
	h.RoomsMu.RLock()
	room, exists := h.Rooms[roomName]
	h.RoomsMu.RUnlock()
	if !exists {
		return
	}
	room.RoomMu.Lock()
	defer room.RoomMu.Unlock()
	for client := range room.Clients {
		if client == sender {
			continue
		}
		client.SendMessage(data)
	}
}
