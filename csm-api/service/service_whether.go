package service

import (
	"csm-api/api"
	"csm-api/config"
	"csm-api/entity"
	"encoding/json"
	"fmt"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-03-05
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 기상청 날씨 api 조회
type ServiceWhether struct {
	ApiKey *config.ApiConfig
}

// func: 기상청 초단기실황 api 조회
// @param
// - date string: 현재날짜(mmdd), time string: 현재시간(hhmm), nx int: 위도변환값, ny int: 경도변환값
func (s *ServiceWhether) GetWhetherSrtNcst(date string, time string, nx int, ny int) (entity.WhetherSrtEntityRes, error) {
	// 초단기실황 url
	url := fmt.Sprintf("http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getUltraSrtNcst?dataType=JSON&ServiceKey=%s&base_date=%s&base_time=%s&nx=%d&ny=%d",
		s.ApiKey.WhetherApiKey,
		date,
		time,
		nx,
		ny,
	)

	// api call
	body, err := api.CallGetAPI(url)
	if err != nil {
		return nil, fmt.Errorf("call GetWhetherSrtNcst API error: %v", err)
	}

	// api response item struct
	type whether struct {
		Response struct {
			Body struct {
				Items entity.WhetherSrtItems `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	// response parse
	var res whether
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, fmt.Errorf("WhetherSrt api JSON parse err: %v", err)
	}

	// whether api response -> go api response convert
	items := res.Response.Body.Items
	whetherRes := entity.WhetherSrtEntityRes{}
	for _, item := range items.Item {
		temp := entity.WhetherSrtEntity{}
		temp.Key = item.Category
		if item.Category == "VEC" {
			temp.Value = entity.WhetherVecString(item.ObsrValue)
		} else {
			temp.Value = item.ObsrValue
		}
		whetherRes = append(whetherRes, temp)
	}

	return whetherRes, nil
}
