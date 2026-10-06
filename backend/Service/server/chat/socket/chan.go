package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// writePump 写协程：单 goroutine 消费 Send 通道
func (c *Client) writePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			// 通道已关闭
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			// 可选优化：合并通道中积压的消息，减少系统调用
			n := len(c.Send)
			for i := 0; i < n; i++ {
				extraMsg := <-c.Send
				c.Conn.WriteMessage(websocket.TextMessage, extraMsg)
			}
		case <-ticker.C:
			// 定时 Ping
			c.Conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage 向客户端发送消息（非阻塞）
func (c *Client) SendMessage(data []byte) {
	select {
	case c.Send <- data:
	// 成功投递
	default:
		// 通道满，丢弃消息并关闭连接（慢客户端保护）
		log.Printf("[%s] 写通道已满，关闭连接", c.Account)
		close(c.Send)
	}
}
