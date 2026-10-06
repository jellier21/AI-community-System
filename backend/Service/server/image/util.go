package image

import (
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"os"
)

type Response struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

// 每次请求时才读环境变量，避免包初始化时 .env 还没加载
func getConfig() (ak, sk, bucket, domain string, mac *credentials.Credentials) {
	ak = os.Getenv("QINIU_ACCESS_KEY")
	sk = os.Getenv("QINIU_SECRET_KEY")
	bucket = os.Getenv("QINIU_BUCKET")
	domain = os.Getenv("QINIU_DOMAIN")
	mac = credentials.NewCredentials(ak, sk)
	return
}
