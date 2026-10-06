package sqlTools

import (
	"database/sql"
	"log"
)

// 本周公告
type WeeklyNotice struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"create_at"`
}

// 本周医疗咨询
type WeeklyMedical struct {
	Symptom   string `json:"symptom"`
	Diagnosis string `json:"diagnosis"`
	CreatedAt string `json:"create_at"`
}

// 查某社区最近 7 天公告
func GetWeeklyNotices(communityUid string, db *sql.DB) ([]WeeklyNotice, error) {
	rows, err := db.Query(
		"SELECT title, content, create_at FROM notice WHERE uid = ? AND create_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) ORDER BY create_at DESC",
		communityUid)
	if err != nil {
		log.Println("查询本周公告失败:", err)
		return nil, err
	}
	defer rows.Close()
	var list []WeeklyNotice
	for rows.Next() {
		var n WeeklyNotice
		if err := rows.Scan(&n.Title, &n.Content, &n.CreatedAt); err != nil {
			continue
		}
		list = append(list, n)
	}
	return list, nil
}

// 查某社区最近 7 天医疗咨询（medical 表无社区号，通过 join 表转成员账号再过滤）
func GetWeeklyMedical(communityUid string, db *sql.DB) ([]WeeklyMedical, error) {
	rows, err := db.Query(
		"SELECT m.symptom, m.diagnosis, m.create_at FROM medical m "+
			"WHERE m.create_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) "+
			"AND m.account IN (SELECT account FROM `join` WHERE community = ? AND is_join = 1) "+
			"ORDER BY m.create_at DESC",
		communityUid)
	if err != nil {
		log.Println("查询本周医疗记录失败:", err)
		return nil, err
	}
	defer rows.Close()
	var list []WeeklyMedical
	for rows.Next() {
		var m WeeklyMedical
		if err := rows.Scan(&m.Symptom, &m.Diagnosis, &m.CreatedAt); err != nil {
			continue
		}
		list = append(list, m)
	}
	return list, nil
}
