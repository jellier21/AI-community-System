package mapreading

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"mime/multipart"
	"net/http"
	"os"
)

func MapReadingClient(account string, file *multipart.FileHeader) (result map[string]interface{}, err error) {
	client := resty.New()
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	resp, err := client.R().
		SetFileReader("file", file.Filename, f).
		SetFormData(map[string]string{"account": account}).
		Post(fmt.Sprintf("%s/api/map/recognize", os.Getenv("CLIENT_URL")))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("下游 /api/map/recognize 返回 %d: %s", resp.StatusCode(), resp.String())
	}
	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	return data, nil
}
