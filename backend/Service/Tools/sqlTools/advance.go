package sqlTools

import (
	"database/sql"
	"log"
)

type Advance struct {
	Id           int    `json:"id"`
	Account      string `json:"account"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CommunityUid string `json:"community_uid"`
}

// 用户提交建议
func SubmitAdvance(account, communityUid, title, content string, db *sql.DB) error {
	sqlStr := "INSERT INTO advance (account, title, content, community_uid) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, account, title, content, communityUid)
	if err != nil {
		log.Println("提交建议失败:", err)
		return err
	}
	return nil
}

// leader 查看本社区所有建议
func GetAdvanceList(communityUid string, db *sql.DB) ([]Advance, error) {
	sqlStr := "SELECT id, account, title, content, community_uid FROM advance WHERE community_uid = ? ORDER BY id DESC"
	rows, err := db.Query(sqlStr, communityUid)
	if err != nil {
		log.Println("查询建议列表失败:", err)
		return nil, err
	}
	defer rows.Close()

	var list []Advance
	for rows.Next() {
		var a Advance
		err := rows.Scan(&a.Id, &a.Account, &a.Title, &a.Content, &a.CommunityUid)
		if err != nil {
			log.Println("扫描建议失败:", err)
			continue
		}
		list = append(list, a)
	}
	return list, nil
}
