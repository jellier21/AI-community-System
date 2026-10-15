package cryptography

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// =====================密码加密保护============================
func PasswordToHash(password string) (string, error) {
	// 密码===哈希===>bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Println("密码哈希错误:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func HashToPasswordBool(hashString string, inputWord string) error {
	//核对密码是否正确
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashString), //待比较的哈希
		[]byte(inputWord),  //输入的密码
	)
	if err != nil {
		fmt.Println("密码不一致")
		return err
	}
	return nil
}
