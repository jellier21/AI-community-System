package util

import (
	"encoding/base64"
	"io"
	"mime/multipart"
)

func ImagefromFormBase64(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	return b64, nil
}
