package service

import (
	"csm-api/api"
	"csm-api/config"
	"csm-api/entity"
	"csm-api/utils"
	"encoding/json"
	"fmt"
)

type ServiceRestDate struct {
	ApiKey *config.ApiConfig
}

// func: 공휴일 날짜 조회 api
// @param
// -
func (s *ServiceRestDate) GetRestDelDates(year string, month string) (entity.RestDates, error) {
	url := fmt.Sprintf("http://apis.data.go.kr/B090041/openapi/service/SpcdeInfoService/getRestDeInfo?_type=json&solYear=%s&solMonth=%s&numOfRows=100&ServiceKey=%s",
		year,
		month,
		s.ApiKey.DataGoApiKey,
	)

	body, err := api.CallGetAPI(url)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	type RestDelInfo struct {
		Response struct {
			Body struct {
				Items entity.RestDels `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	var res RestDelInfo
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		fmt.Println("Unmarshal err or non RestDel Info")
		return entity.RestDates{}, nil
	}

	items := res.Response.Body.Items
	resRests := entity.RestDates{}
	for _, item := range items.Item {
		rest := entity.RestDate{}
		rest.Reason = item.DateName
		rest.RestDate = item.Locdate
		resRests = append(resRests, rest)
	}

	return resRests, nil
}
