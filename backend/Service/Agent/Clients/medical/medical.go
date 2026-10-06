package medical

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-resty/resty/v2"
	"mime/multipart"
	"net/http"
	"os"
)

func TestClient(msg string) (result string, err error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("message", msg).
		Get(fmt.Sprintf("%s/", os.Getenv("CLIENT_URL")))
	if err != nil {
		return "", err
	}
	return response.String(), nil
}

func TextAgentClient(msg string) (result map[string]interface{}, err error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("message", msg).
		Get(fmt.Sprintf("%s/api/text", os.Getenv("CLIENT_URL")))
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func MedicalAgentClient(account, question string, file *multipart.FileHeader) (result map[string]interface{}, err error) {
	client := resty.New()
	var resp *resty.Response
	if file != nil {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		resp, err = client.R().
			SetFileReader("file", file.Filename, f).
			SetFormData(map[string]string{"account": account, "ask": question}).
			Post(fmt.Sprintf("%s/api/medical/cure", os.Getenv("CLIENT_URL")))
	} else {
		resp, err = client.R().
			SetFormData(map[string]string{"account": account, "ask": question}).
			Post(fmt.Sprintf("%s/api/medical/cure", os.Getenv("CLIENT_URL")))
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("下游 /api/medical/cure 返回 %d: %s", resp.StatusCode(), resp.String())
	}
	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	return data, nil
}
