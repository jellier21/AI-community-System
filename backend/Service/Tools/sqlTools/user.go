package sqlTools

import (
	"Service/Tools/JWT"
	"Service/Tools/cryptography"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func RegisterTools(_account_, _password_ string, db *sql.DB) error {
	// 先查重
	var existCnt int
	if err := db.QueryRow("SELECT COUNT(*) FROM user WHERE account = ?", _account_).Scan(&existCnt); err != nil {
		log.Println("查重失败:", err)
		return err
	}
	if existCnt > 0 {
		return errors.New("账号已存在")
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	tempNick := fmt.Sprintf("user_%05d", count+1)
	tempPassword, err := cryptography.PasswordToHash(_password_)
	if err != nil {
		log.Println("密码哈希失败:", err)
		return err
	}

	sqlStr := "INSERT INTO user (account, password, nick) VALUES (?, ?, ?)"
	_, err = db.Exec(sqlStr, _account_, tempPassword, tempNick)
	if err != nil {
		log.Println("注册用户失败:", err)
	}
	return err
}

func LoginTools(_account_, _password_ string, db *sql.DB) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM user WHERE account = ?", _account_).Scan(&password)
	if err != nil {
		log.Println("账号错误")
		return "", err
	}
	if err := cryptography.HashToPasswordBool(password, _password_); err != nil {
		log.Println("密码错误:", err)
		return "", fmt.Errorf("密码错误")
	}
	token, err := JWT.GenerateJsonWebToken(_account_)
	if err != nil {
		log.Println("生成Token失败:", err)
		return "", err
	}
	return token, nil
}

func QueryUserNickname(_account_ string, db *sql.DB) (nick string, err error) {
	sqlStr := "select nick from user where account = ?"
	err = db.QueryRow(sqlStr, _account_).Scan(&nick)
	return nick, err
}

func AmendNickByAccount(_nick_, _account_ string, db *sql.DB) error {
	sqlStr := "update user set nick = ? where account = ?"
	_, err := db.Exec(sqlStr, _nick_, _account_)
	return err
}
func AmendPasswordByAccount(_account_, _new_password_, _old_password_ string, db *sql.DB) error {
	var PasswordDB string
	err := db.QueryRow("SELECT password FROM user WHERE account = ?", _account_).Scan(&PasswordDB)
	if err != nil {
		log.Println("账号错误")
		return err
	}

	if err = cryptography.HashToPasswordBool(PasswordDB, _old_password_); err != nil {
		log.Println("密码错误:", err)
		return fmt.Errorf("原密码错误")
	}
	newPassword, err := cryptography.PasswordToHash(_new_password_)
	if err != nil {
		log.Println("新密码哈希失败:", err)
		return err
	}
	sqlStr := "update user set password = ? where account = ?"
	_, err = db.Exec(sqlStr, newPassword, _account_)
	return err
}

func AmendMySignatureByAccount(_signature_, _account_ string, db *sql.DB) error {
	sqlStr := "update user set signature = ? where account = ?"
	_, err := db.Exec(sqlStr, _signature_, _account_)
	return err
}

// 查头像
func QueryUserAvatar(account string, db *sql.DB) (string, error) {
	var avatar string
	err := db.QueryRow("SELECT avatar FROM user WHERE account = ?", account).Scan(&avatar)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return avatar, err
}

// 更新头像
func AmendAvatarByAccount(avatar, account string, db *sql.DB) error {
	_, err := db.Exec("UPDATE user SET avatar = ? WHERE account = ?", avatar, account)
	return err
}

type UserProfile struct {
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

func QueryUserProfile(account string, db *sql.DB) (UserProfile, error) {
	var p UserProfile
	p.Account = account
	err := db.QueryRow("SELECT IFNULL(nick,''), IFNULL(avatar,''), IFNULL(signature,'') FROM user WHERE account = ?", account).
		Scan(&p.Nickname, &p.Avatar, &p.Signature)
	return p, err
}
