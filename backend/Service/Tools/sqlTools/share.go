package sqlTools

import (
	"database/sql"
	"errors"
	"log"
)

type Share struct {
	Id           int    `json:"id"`
	Account      string `json:"account"`
	CommunityUid string `json:"community_uid"`
	Content      string `json:"content"`
	Images       string `json:"images"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"` // ← 加这行
	CreateAt     string `json:"create_at"`
	UpdateAt     string `json:"update_at"`
}

// 发布动态
func PublishShare(account, communityUid, content, images string, db *sql.DB) error {
	sqlStr := "INSERT INTO share (account, community_uid, content, images) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, account, communityUid, content, images)
	if err != nil {
		log.Println("发布动态失败:", err)
		return err
	}
	return nil
}

// 查我自己的动态
func GetMyShares(account string, db *sql.DB) ([]Share, error) {
	sqlStr := "SELECT s.id, s.account, s.community_uid, s.content, s.images, s.create_at, s.update_at, IFNULL(j.name, s.account), IFNULL(u.avatar, '') FROM share s LEFT JOIN `join` j ON s.account = j.account AND s.community_uid = j.community LEFT JOIN user u ON s.account = u.account WHERE s.community_uid = ? ORDER BY s.id DESC"
	rows, err := db.Query(sqlStr, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShares(rows)
}

// 查本社区所有动态（管理员用）
func GetCommunityShares(communityUid string, db *sql.DB) ([]Share, error) {
	sqlStr := "SELECT s.id, s.account, s.community_uid, s.content, s.images, s.create_at, s.update_at, IFNULL(j.name, s.account), IFNULL(u.avatar, '') FROM share s LEFT JOIN `join` j ON s.account = j.account AND s.community_uid = j.community LEFT JOIN user u ON s.account = u.account WHERE s.community_uid = ? ORDER BY s.id DESC"
	rows, err := db.Query(sqlStr, communityUid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShares(rows)
}

// 查全站最新动态（首页推荐位用，不限制社区）
func GetLatestShares(limit int, db *sql.DB) ([]Share, error) {
	sqlStr := "SELECT s.id, s.account, s.community_uid, s.content, s.images, s.create_at, s.update_at, IFNULL(u.nick, s.account), IFNULL(u.avatar, '') FROM share s LEFT JOIN user u ON s.account = u.account ORDER BY s.id DESC LIMIT ?"
	rows, err := db.Query(sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShares(rows)
}

// 删除动态（校验：本人 或 是该社区 leader）
func DeleteShare(shareId int, account string, db *sql.DB) error {
	// 先查这条动态
	var s Share
	err := db.QueryRow("SELECT id, account, community_uid FROM share WHERE id = ?", shareId).
		Scan(&s.Id, &s.Account, &s.CommunityUid)
	if err != nil {
		return err
	}
	// 本人发的可以删
	if s.Account == account {
		_, err = db.Exec("DELETE FROM share WHERE id = ?", shareId)
		return err
	}
	// 不是本人，检查是不是 leader
	err = CheckLeaderOfCommunity(account, s.CommunityUid, db)
	if err == nil {
		_, err = db.Exec("DELETE FROM share WHERE id = ?", shareId)
		return err
	}
	return err
}

// 修改动态（只能改自己的）
func UpdateShare(shareId int, account, content, images string, db *sql.DB) error {
	res, err := db.Exec("UPDATE share SET content = ?, images = ? WHERE id = ? AND account = ?",
		content, images, shareId, account)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("无权修改或动态不存在")
	}
	return nil
}

func scanShares(rows *sql.Rows) ([]Share, error) {
	var list []Share
	for rows.Next() {
		var s Share
		err := rows.Scan(&s.Id, &s.Account, &s.CommunityUid, &s.Content, &s.Images,
			&s.CreateAt, &s.UpdateAt, &s.Name, &s.Avatar)

		if err != nil {
			log.Println("扫描动态失败:", err)
			continue
		}
		list = append(list, s)
	}
	return list, nil
}
