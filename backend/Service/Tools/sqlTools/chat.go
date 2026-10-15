package sqlTools

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type MessageDB struct {
	ID        int64  `json:"id"`
	RoomUid   string `json:"room_uid"`
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

func GetRoomMessages(roomUid string, page, pageSize int, db *sql.DB) ([]MessageDB, error) {
	offset := (page - 1) * pageSize
	sqlStr := `SELECT id, room, account, content, time 
               FROM message 
               WHERE room = ? 
               ORDER BY time DESC 
               LIMIT ? OFFSET ?`

	rows, err := db.Query(sqlStr, roomUid, pageSize, offset)
	if err != nil {
		log.Println("查询消息历史失败:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []MessageDB
	for rows.Next() {
		var msg MessageDB
		var createdAt sql.NullString // 使用可空字符串
		err := rows.Scan(&msg.ID, &msg.RoomUid, &msg.Account,
			&msg.Content, &createdAt)
		msg.Nickname = GetMemberDisplayName(roomUid, msg.Account, db)
		msg.Type = "message"
		if err != nil {
			log.Println("扫描消息失败:", err)
			continue
		}
		if createdAt.Valid {
			msg.CreatedAt = createdAt.String
		}
		messages = append(messages, msg)
	}

	// 反转消息顺序（前端显示需要按时间正序）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func GetRecentMessages(roomUid string, limit int, db *sql.DB) ([]MessageDB, error) {
	sqlStr := `SELECT m.id, m.room, m.account, m.content, m.time,
		IFNULL(u.avatar, '')
		FROM message m
		LEFT JOIN user u ON m.account = u.account
		WHERE m.room = ?
		ORDER BY m.time DESC
		LIMIT ?`

	rows, err := db.Query(sqlStr, roomUid, limit)
	if err != nil {
		log.Println("查询最近消息失败:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []MessageDB
	for rows.Next() {
		var msg MessageDB
		var createdAt sql.NullString
		err := rows.Scan(&msg.ID, &msg.RoomUid, &msg.Account,
			&msg.Content, &createdAt, &msg.Avatar)
		msg.Nickname = GetMemberDisplayName(roomUid, msg.Account, db)
		msg.Type = "message"
		if err != nil {
			log.Println("扫描消息失败:", err)
			continue
		}
		if createdAt.Valid {
			msg.CreatedAt = createdAt.String
		}
		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ============================ 房间SQL    ===========================
type RoomDB struct {
	RoomUid   string `json:"room_uid"`
	RoomName  string `json:"room_name"`
	Leader    string `json:"leader"`
	CreatedAt string `json:"created_at"`
}

func GetChatroomLists(account string, db *sql.DB) ([]RoomDB, error) {
	// 私聊房间的 r.name 是创建者当时看到的对方昵称，对被加方来说会显示成自己。
	// 因此对 type='private' 的房间，按当前登录者动态取"另一个成员"的昵称作为显示名；
	// 群聊/社区房继续用 r.name。
	sqlStr := `SELECT r.uid,
		CASE WHEN r.type = 'private' THEN
			(SELECT IFNULL(u2.nick, g2.account)
			 FROM ` + "`groups`" + ` g2
			 LEFT JOIN user u2 ON g2.account = u2.account
			 WHERE g2.uid = r.uid AND g2.account != ?
			 LIMIT 1)
		ELSE r.name END AS display_name,
		r.leader, r.updatetime
		FROM ` + "`groups`" + ` g
		JOIN room r ON g.uid = r.uid
		WHERE g.account = ?`
	rows, err := db.Query(sqlStr, account, account)
	if err != nil {
		log.Println("查询房间列表失败:", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []RoomDB
	for rows.Next() {
		var room RoomDB
		var updatetime sql.NullString
		err := rows.Scan(&room.RoomUid, &room.RoomName, &room.Leader, &updatetime)
		if err != nil {
			log.Println("扫描房间失败:", err)
			continue
		}
		if updatetime.Valid {
			room.CreatedAt = updatetime.String
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// SaveMessage 保存消息到数据库
func SaveMessage(roomUid, account, content string, db *sql.DB) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := `INSERT INTO message (room, account, content, time) 
               VALUES (?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, roomUid, account, content, currentTime)
	if err != nil {
		log.Println("保存消息失败:", err)
		return err
	}
	return nil
}

func GetUserAvatarByAccount(account string, db *sql.DB) string {
	var avatar string
	err := db.QueryRow("SELECT IFNULL(avatar,'') FROM user WHERE account = ?", account).Scan(&avatar)
	if err != nil {
		return ""
	}
	return avatar
}

// ============================ 私聊 / 群聊 / 好友 ===========================

// IsUserExists 查用户是否存在，返回昵称
func IsUserExists(account string, db *sql.DB) (bool, string) {
	var nick string
	err := db.QueryRow("SELECT IFNULL(nick,'') FROM user WHERE account = ?", account).Scan(&nick)
	if err != nil {
		return false, ""
	}
	return true, nick
}

// PrivateRoomUid 生成稳定的私聊房间号（双方顺序固定）
func PrivateRoomUid(a, b string) string {
	if a < b {
		return "p_" + a + "_" + b
	}
	return "p_" + b + "_" + a
}

// CreatePrivateRoom 创建私聊房间，myAccount=发起者，friendAccount=对方账号，friendNick=对方昵称
func CreatePrivateRoom(myAccount, friendAccount, friendNick string, db *sql.DB) error {
	roomUid := PrivateRoomUid(myAccount, friendAccount)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO room (uid, name, leader, type) VALUES (?, ?, ?, 'private')",
		roomUid, friendNick, myAccount)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("INSERT INTO `groups` (account, uid) VALUES (?,?), (?,?)",
		myAccount, roomUid, friendAccount, roomUid)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// CreateGroupRoom 创建群聊
func CreateGroupRoom(myAccount, groupName string, members []string, db *sql.DB) (string, error) {
	roomUid := "g_" + myAccount + "_" + fmt.Sprintf("%d", time.Now().UnixMilli())
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	_, err = tx.Exec("INSERT INTO room (uid, name, leader, type) VALUES (?, ?, ?, 'group')",
		roomUid, groupName, myAccount)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	_, err = tx.Exec("INSERT INTO `groups` (account, uid) VALUES (?, ?)", myAccount, roomUid)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	for _, m := range members {
		if m == "" || m == myAccount {
			continue
		}
		_, err = tx.Exec("INSERT INTO `groups` (account, uid) VALUES (?, ?)", m, roomUid)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	return roomUid, tx.Commit()
}

// FriendItem 好友列表项
type FriendItem struct {
	RoomUid       string `json:"room_uid"`
	FriendAccount string `json:"friend_account"`
	FriendNick    string `json:"friend_nick"`
	FriendAvatar  string `json:"friend_avatar"`
}

// GetFriendList 查我的所有私聊好友
func GetFriendList(myAccount string, db *sql.DB) ([]FriendItem, error) {
	sqlStr := `SELECT r.uid,
		(SELECT g2.account FROM ` + "`groups`" + ` g2 WHERE g2.uid = r.uid AND g2.account != ? LIMIT 1),
		IFNULL(u.nick,''), IFNULL(u.avatar,'')
		FROM room r
		JOIN ` + "`groups`" + ` g ON r.uid = g.uid
		LEFT JOIN user u ON u.account = (SELECT g2.account FROM ` + "`groups`" + ` g2 WHERE g2.uid = r.uid AND g2.account != ? LIMIT 1)
		WHERE g.account = ? AND r.type = 'private'`
	rows, err := db.Query(sqlStr, myAccount, myAccount, myAccount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []FriendItem
	for rows.Next() {
		var f FriendItem
		if err := rows.Scan(&f.RoomUid, &f.FriendAccount, &f.FriendNick, &f.FriendAvatar); err != nil {
			continue
		}
		list = append(list, f)
	}
	return list, nil
}

// DeletePrivateRoom 删除私聊房间及相关数据
func DeletePrivateRoom(roomUid string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	tx.Exec("DELETE FROM room WHERE uid = ? AND type = 'private'", roomUid)
	tx.Exec("DELETE FROM `groups` WHERE uid = ?", roomUid)
	tx.Exec("DELETE FROM message WHERE room = ?", roomUid)
	return tx.Commit()
}
