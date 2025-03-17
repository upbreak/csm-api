package service

import (
	"csm-api/api"
	"csm-api/config"
	"csm-api/entity"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

/**
 * @author 작성자: 정지영
 * @created 작성일: 2025-03-13
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 브이월드 주소 api 조회
type ServiceAddressSearching struct {
	ApiKey *config.ApiConfig
}

// func: 도로명주소로 위도, 경도 조회
// @param
// - roadAddress string: 도로명주소
func (s *ServiceAddressSearching) GetAPILatitudeLongtitude(roadAddress string) (*entity.Point, error) {

	if roadAddress == "" {
		return nil, fmt.Errorf("roadAddress parameter is missing")
	}
	// apiUrl주소
	apiUrl := fmt.Sprintf(`https://api.vworld.kr/req/search?key=%s&version=%s&service=%s&request=%s&type=%s&crs=%s&category=%s&query=%s&domain=%s&format=%s&errorFormat=%s`,
		url.QueryEscape(s.ApiKey.VworldApiKey),
		"2.0",
		"search",
		"search",
		"address",
		"EPSG:4326",
		"ROAD",
		url.QueryEscape(roadAddress),
		"csm.htenc.co.kr",
		"json",
		"json",
	)

	// api 호출
	body, err := api.CallGetAPI(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("call GetWhetherSrtNcst API error: %v", err)
	}

	// api response struct
	type AddressItem struct {
		Id      string   `json:"id"`
		Address struct{} `json:"address"`
		Point   struct {
			Longitude string `json:"x"`
			Latitude  string `json:"y"`
		} `json:"point"`
	}
	type AddressItems []AddressItem

	type AddressSearching struct {
		Response struct {
			Status string `json:"status"`
			Result struct {
				Crs   string       `json:"crs"`
				Type  string       `json:"type"`
				Items AddressItems `json:"items"`
			} `json:"result"`
			Error struct {
				Level string `json:"level"`
				Code  int    `json:"code"`
				Text  string `json:"text"`
			}
		} `json:"response"`
	}

	// 반환 경/위도 객체
	point := &entity.Point{}

	// api response parse
	var res AddressSearching
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return point, fmt.Errorf("AddressSearch api JSON parse err: %v", err)
	}

	if res.Response.Status == "ERROR" {
		return point, fmt.Errorf("AddressSearch api response err: %s", res.Response.Error.Text)

	} else {
		items := res.Response.Result.Items

		// 반환 값 추가
		point.Latitude, _ = strconv.ParseFloat(items[0].Point.Latitude, 64)
		point.Longitude, _ = strconv.ParseFloat(items[0].Point.Longitude, 64)
	}

	return point, nil
}
