package sqlTools

import (
	"database/sql"
	"log"
)

type Notice struct {
	Id        int    `json:"id"`
	Uid       string `json:"uid"` // 社区号
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"create_at"`
}

// 发布公告（uid 就是社区号）
func PublishNotice(communityUid, title, content string, db *sql.DB) error {
	sqlStr := "INSERT INTO notice (uid, title, content, create_at) VALUES (?, ?, ?, NOW())"
	_, err := db.Exec(sqlStr, communityUid, title, content)
	if err != nil {
		log.Println("发布公告失败:", err)
		return err
	}
	return nil
}

// 查某社区所有公告（按时间倒序）
func GetNoticeList(communityUid string, db *sql.DB) ([]Notice, error) {
	sqlStr := "SELECT id, uid, title, content, create_at FROM notice WHERE uid = ? ORDER BY create_at DESC"
	rows, err := db.Query(sqlStr, communityUid)
	if err != nil {
		log.Println("查询公告列表失败:", err)
		return nil, err
	}
	defer rows.Close()

	var list []Notice
	for rows.Next() {
		var n Notice
		err := rows.Scan(&n.Id, &n.Uid, &n.Title, &n.Content, &n.CreatedAt)
		if err != nil {
			log.Println("扫描公告失败:", err)
			continue
		}
		list = append(list, n)
	}
	return list, nil
}

// 查某社区最新一条
func GetLatestNotice(communityUid string, db *sql.DB) (Notice, error) {
	var n Notice
	sqlStr := "SELECT id, uid, title, content, create_at FROM notice WHERE uid = ? ORDER BY create_at DESC LIMIT 1"
	err := db.QueryRow(sqlStr, communityUid).Scan(&n.Id, &n.Uid, &n.Title, &n.Content, &n.CreatedAt)
	if err == sql.ErrNoRows {
		return Notice{}, nil
	}
	return n, err
}
