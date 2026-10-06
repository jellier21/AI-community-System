package JWT

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

func GenerateJsonWebToken(uid string) (string, error) {
	claims := JWTStruct{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 优先从 Authorization 头获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = authHeader
		} else {
			// WebSocket 连接从查询参数获取
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			log.Println("未提供认证令牌")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		//log.Println("token使用:", tokenString)
		token, err := jwt.ParseWithClaims(tokenString, &JWTStruct{}, func(token *jwt.Token) (interface{}, error) {
			return GetJWTSecret(), nil
		})
		if err != nil {
			log.Println("令牌已过期:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌无效"})
			c.Abort()
			return
		}
		c.Set("userID", token.Claims.(*JWTStruct).UserID)
		c.Next()
	}
}

func JwtExpirationMonitoring(r *gin.Engine) {
	r.GET("/api/jwt/monitor", AuthMiddleware(), func(c *gin.Context) {
		log.Println("JWT 鉴权检测 已通过")
		c.JSON(http.StatusOK, gin.H{"message": "JWT 鉴权检测 已通过"})
	})
}
