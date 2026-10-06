package JWT

import (
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	jwtSecretOnce sync.Once
)

func GetJWTSecret() []byte {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	return jwtSecret
}

type JWTStruct struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}
