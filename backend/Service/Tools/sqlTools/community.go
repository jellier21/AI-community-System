package sqlTools

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type Community struct {
	Uid         string `json:"uid"`
	Leader      string `json:"leader"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CheckWheatherIs(account string, db *sql.DB) error {
	sqlStr := "SELECT * FROM community WHERE leader = ?"
	rows, err := db.Query(sqlStr, account)
	if rows != nil && rows.Next() {
		log.Println("该用户已存在社区")
		return errors.New("该用户已存在社区")
	}
	if err != nil {
		log.Println("查询失败:", err)
		return err
	}
	return err
}

func BuildCommunitySQL(_account_, communityName, communityDescription, _uid_ string, db *sql.DB) error {
	err := CheckWheatherIs(_account_, db)
	if err != nil {
		return err
	}
	// 1. 建社区
	sqlStr := "INSERT INTO community (uid, leader, name, description) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, _uid_, _account_, communityName, communityDescription)
	if err != nil {
		return err
	}
	// 2. leader 自己加入 join 表
	sqlStr = "INSERT INTO `join` (account,community,is_join,name) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, _account_, _uid_, 1, "负责人")
	if err != nil {
		return err
	}
	// 3. 创建同名聊天室（房间号=社区号，房间名=社区名）
	sqlStr = "INSERT INTO room (uid, name, leader, updatetime) VALUES (?, ?, ?, NOW())"

	_, err = db.Exec(sqlStr, _uid_, communityName, _account_)
	if err != nil {
		log.Println("创建聊天室失败:", err)
		return err
	}
	// 4. leader 加入聊天室
	sqlStr = "INSERT INTO `groups` (uid, account) VALUES (?, ?)"
	_, err = db.Exec(sqlStr, _uid_, _account_)
	if err != nil {
		log.Println("leader 加入聊天室失败:", err)
		return err
	}
	return nil
}

type HavingCommunity struct {
	Type                 string `json:"type"`
	Uid                  string `json:"uid"`
	CommunityName        string `json:"community_name"`
	CommunityDescription string `json:"community_description"`
	Leader               string `json:"leader"`
}

func UserWheatherHaveCommunity(account string, db *sql.DB) (HC HavingCommunity, err error) {
	sqlStr := "SELECT uid,name,description,leader FROM community WHERE leader = ?"
	err = db.QueryRow(sqlStr, account).Scan(&HC.Uid, &HC.CommunityName, &HC.CommunityDescription, &HC.Leader)
	HC.Type = "leader"
	if err != nil {
		log.Println("未查询到用户名下有社区，现在查询是否加入社区")
		sqlStr = "SELECT community FROM `join` WHERE account = ? AND is_join = 1"

		err = db.QueryRow(sqlStr, account).Scan(&HC.Uid)
		if err != nil {
			log.Println("未查询到用户名下有社区")
			return HavingCommunity{}, errors.New("未查询到用户加入社区")
		}

		var TempCommunity Community
		TempCommunity, err = GetCommunityInfo(HC.Uid, db)
		HC.CommunityName = TempCommunity.Name
		HC.CommunityDescription = TempCommunity.Description
		HC.Leader = TempCommunity.Leader
		HC.Type = "join"
	}
	return HC, nil
}

func GetCommunityInfo(Uid string, db *sql.DB) (Community, error) {
	sqlStr := "SELECT uid,name,description,leader FROM community WHERE uid = ?"
	rows, err := db.Query(sqlStr, Uid)
	if err != nil {
		log.Println("查询失败:", err)
		return Community{}, err
	}
	if rows.Next() {
		var community Community
		err = rows.Scan(&community.Uid, &community.Leader, &community.Name, &community.Description)
		if err != nil {
			log.Println("扫描失败:", err)
			return Community{}, err
		}
		return community, nil
	}
	return Community{}, nil
}

func QuitCommunity(account string, db *sql.DB) error {
	log.Println("account:", account)
	sqlStr := "SELECT leader FROM community WHERE leader = ?"
	var CommunityUid string
	err := db.QueryRow(sqlStr, account).Scan(&CommunityUid)
	if err != nil {
		sqlStr = "SELECT community FROM `join` WHERE account = ?"
		err = db.QueryRow(sqlStr, account).Scan(&CommunityUid)
		if err != nil {
			log.Println("未查询到用户加入社区")
			return errors.New("未查询到用户加入社区")
		}
		sqlDel := "DELETE FROM `join` WHERE account = ?"
		_, err = db.Exec(sqlDel, account)
		if err != nil {
			log.Println("退出社区失败:", err)
			return err
		}
		// 同步移出聊天室
		sqlDel = "DELETE FROM `groups` WHERE account = ?"
		_, err = db.Exec(sqlDel, account)
		if err != nil {
			log.Println("移出聊天室失败:", err)
			return err
		}
		return nil

	}
	sqlDel := "DELETE FROM community WHERE leader = ?"
	_, err = db.Exec(sqlDel, CommunityUid)
	if err != nil {
		log.Println("解散社区失败:", err)
		return err
	}
	sqlDel = "DELETE FROM `join` WHERE account = ?"
	_, err = db.Exec(sqlDel, account)
	if err != nil {
		log.Println("退出社区失败:", err)
		return err
	}
	// 同步删除聊天室和所有群成员
	sqlDel = "DELETE FROM room WHERE uid = ?"
	_, err = db.Exec(sqlDel, CommunityUid)
	if err != nil {
		log.Println("删除聊天室失败:", err)
		return err
	}
	sqlDel = "DELETE FROM `groups` WHERE uid = ?"
	_, err = db.Exec(sqlDel, CommunityUid)
	if err != nil {
		log.Println("清空群成员失败:", err)
		return err
	}
	return nil

}

func UserWheatherIsLeader(_account_ string, db *sql.DB) error {
	var name string
	sqlStr := "SELECT name FROM community WHERE leader = ?"
	err := db.QueryRow(sqlStr, _account_).Scan(&name)
	if err != nil {
		log.Println("未查询到用户为社区负责人")
		return errors.New("未查询到用户为社区领袖")
	}
	return nil
}

func GetCommunityUID(account string, db *sql.DB) (string, error) {
	var uid string
	// 先查是不是 leader
	err := db.QueryRow("SELECT uid FROM community WHERE leader = ?", account).Scan(&uid)
	if err == nil {
		return uid, nil
	}
	// 不是 leader，查 join 表（已加入的成员）
	err = db.QueryRow("SELECT community FROM `join` WHERE account = ? AND is_join = 1", account).Scan(&uid)
	if err != nil {
		log.Println("用户未加入任何社区:", account)
		return "", errors.New("未查询到用户所属社区")
	}
	return uid, nil
}

// ------------------------Member---------------------------
type Member struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Account string `json:"account"`
	IsJoin  int    `json:"is_join"`
}

func GetMyMemberList(Uid string, db *sql.DB) ([]Member, error) {
	sqlStr := "SELECT community,name,account,is_join FROM `join` WHERE community = ?"
	rows, err := db.Query(sqlStr, Uid)
	if err != nil {
		log.Println("查询失败:", err)
		return []Member{}, err
	}
	var members []Member
	for rows.Next() {
		var member Member
		err = rows.Scan(&member.Uid, &member.Name, &member.Account, &member.IsJoin)
		if err != nil {
			log.Println("扫描失败:", err)
			return []Member{}, err
		}
		members = append(members, member)
	}
	return members, nil
}

// delete member-----------------------
// 校验指定账号是不是指定社区的 leader
func CheckLeaderOfCommunity(leaderAccount, communityUid string, db *sql.DB) error {
	var uid string
	sqlStr := "SELECT uid FROM community WHERE leader = ? AND uid = ?"
	err := db.QueryRow(sqlStr, leaderAccount, communityUid).Scan(&uid)
	if err != nil {
		return errors.New("非该社区负责人")
	}
	return nil
}

// 批量移除成员（限定 community，避免误删其他社区）
func RemoveMembers(communityUid string, accounts []string, db *sql.DB) error {
	if len(accounts) == 0 {
		return nil
	}
	placeholders := strings.Repeat("?,", len(accounts))
	placeholders = placeholders[:len(placeholders)-1]

	// 1. 从 join 表删
	sqlStr := "DELETE FROM `join` WHERE community = ? AND account IN (" + placeholders + ")"
	args := make([]interface{}, 0, len(accounts)+1)
	args = append(args, communityUid)
	for _, a := range accounts {
		args = append(args, a)
	}
	_, err := db.Exec(sqlStr, args...)
	if err != nil {
		log.Println("删除成员失败:", err)
		return err
	}

	// 2. 从聊天室 groups 表删
	sqlStr = "DELETE FROM `groups` WHERE uid = ? AND account IN (" + placeholders + ")"
	_, err = db.Exec(sqlStr, args...)
	if err != nil {
		log.Println("移出聊天室失败:", err)
		return err
	}
	return nil
}

// ============== 添加成员 / 加入社区 ==============

// leader 邀请成员：插入 join 表，is_join=0 表示待对方确认
func AddMember(communityUid, targetAccount, memberName string, db *sql.DB) error {
	// 1. 检查对方账号是否存在
	var nick string
	err := db.QueryRow("SELECT nick FROM user WHERE account = ?", targetAccount).Scan(&nick)
	if err != nil {
		log.Println("用户不存在:", targetAccount)
		return errors.New("该账号不存在")
	}

	// 2. 检查是否已经在社区里（包括已邀请未加入的）
	var cnt int
	err = db.QueryRow("SELECT COUNT(*) FROM `join` WHERE community = ? AND account = ?",
		communityUid, targetAccount).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("该用户已在社区中")
	}

	// 3. 插入，name 用 leader 设置的显示名
	sqlStr := "INSERT INTO `join` (account, community, is_join, name) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, targetAccount, communityUid, 0, memberName)
	if err != nil {
		log.Println("添加成员失败:", err)
		return err
	}
	return nil
}

// 成员通过社区号主动加入：把 is_join 从 0 改成 1
func JoinCommunityByUid(communityUid, account string, db *sql.DB) error {
	var cnt int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM `join` WHERE community = ? AND account = ? AND is_join = 0",
		communityUid, account).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("该社区未邀请你，无法加入")
	}

	// 1. 更新 join 表状态
	sqlStr := "UPDATE `join` SET is_join = 1 WHERE community = ? AND account = ? AND is_join = 0"
	_, err = db.Exec(sqlStr, communityUid, account)
	if err != nil {
		log.Println("加入社区失败:", err)
		return err
	}
	// 2. 同步加入聊天室（防重复插入）
	sqlStr = "INSERT IGNORE INTO `groups` (uid, account) VALUES (?, ?)"
	_, err = db.Exec(sqlStr, communityUid, account)
	if err != nil {
		log.Println("加入聊天室失败:", err)
		return err
	}
	return nil
}

// ========================辅助函数==============
// 查房间内某成员的显示名：社区群用 join.name，否则用 user.nick
func GetMemberDisplayName(roomUid, account string, db *sql.DB) string {
	// 先查 join 表（社区群）
	var name string
	err := db.QueryRow("SELECT name FROM `join` WHERE community = ? AND account = ?",
		roomUid, account).Scan(&name)
	if err == nil && name != "" {
		return name
	}
	// 兜底用注册昵称
	nick, err := QueryUserNickname(account, db)
	if err != nil {
		return account
	}
	return nick
}
