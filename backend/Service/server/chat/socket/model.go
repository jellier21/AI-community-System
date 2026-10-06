package socket

import (
	"database/sql"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	// CheckOrigin 返回 true 表示允许所有来源的 WebSocket 连接
	// 生产环境应限制为具体域名
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 包装 WebSocket 连接和用户信息
type Client struct {
	Conn        *websocket.Conn
	CurrentRoom string // 当前所在房间Uid
	Account     string
	Send        chan []byte //写通道-缓存256条
	Hub         *Hub        // Hub 引用
}

const sendBufferSize = 256

// NewClient 创建客户端并启动写协程
func NewClient(conn *websocket.Conn, account string) *Client {
	client := &Client{
		Conn:    conn,
		Account: account,
		Send:    make(chan []byte, sendBufferSize),
	}
	// 启动写协程
	go client.writePump()
	return client
}

var (
	clients = make(map[*Client]bool) // 所有在线连接
	mu      sync.Mutex               // 保护 clients 的互斥锁
)

const (
	// 读超时：如果 60 秒内没收到任何消息（包括 Pong），认为连接已断开
	ReadTimeout = 60 * time.Second
	// 写超时：每次写操作的超时时间
	WriteTimeout = 10 * time.Second
	// Ping 间隔：每 30 秒发一次 Ping
	PingPeriod = 30 * time.Second
)

type Message struct {
	Type     string `json:"type"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content"`
	Time     string `json:"time"`
}

func NewMessage(account, content string) Message {
	//创建一条用户消息
	return Message{
		Type:    "msg",
		Account: account,
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	}
}

func SystemMessage(t, content string) Message {
	// SystemMessage 创建一条系统消息
	return Message{
		Type:    t,  // join / leave / system
		Account: "", // 系统消息无特定用户
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	}
}

// ToJSON 将 Message 序列化为 JSON 字节切片
func (m Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ParseMessage 从 JSON 字节切片解析为 Message
func ParseMessage(data []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// ---------------HUB-------------------------------------------------------------
type Room struct {
	RoomUid string           `json:"roomUid"`
	Clients map[*Client]bool // 房间内所有客户端
	RoomMu  sync.Mutex       // 保护房间信息的互斥锁
}

type Hub struct {
	// 房间管理
	Rooms   map[string]*Room
	RoomsMu sync.RWMutex
	// 调度通道
	Register   chan *ClientRegistration
	Unregister chan *Client
	Broadcast  chan *RoomMessage
	// 数据库连接（用于消息持久化）
	DB *sql.DB
}

// ClientRegistration 注册信息
type ClientRegistration struct {
	Client  *Client
	RoomUid string
}

// RoomMessage 房间消息
type RoomMessage struct {
	RoomUid string
	Data    []byte
	Sender  *Client
}
