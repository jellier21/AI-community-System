package report

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"

	"Service/Tools/sqlTools"
)

type WeeklyReportRequest struct {
	CommunityName string                   `json:"community_name"`
	Notices       []sqlTools.WeeklyNotice  `json:"notices"`
	Medicals      []sqlTools.WeeklyMedical `json:"medicals"`
}

func WeeklyReportClient(communityName string, notices []sqlTools.WeeklyNotice, medicals []sqlTools.WeeklyMedical) (result map[string]interface{}, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(WeeklyReportRequest{
			CommunityName: communityName,
			Notices:       notices,
			Medicals:      medicals,
		}).
		Post(fmt.Sprintf("%s/api/report/generate", os.Getenv("CLIENT_URL")))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("下游 /api/report/generate 返回 %d: %s", resp.StatusCode(), resp.String())
	}
	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	return data, nil
}
