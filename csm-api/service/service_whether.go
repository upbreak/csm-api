package service

import (
	"csm-api/api"
	"csm-api/config"
	"csm-api/entity"
	"encoding/json"
	"fmt"
	"strings"
	"time"
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
	url := fmt.Sprintf("http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getUltraSrtFcst?dataType=JSON&ServiceKey=%s&base_date=%s&base_time=%s&nx=%d&ny=%d&numOfRows=%d",
		s.ApiKey.WhetherApiKey,
		date,
		time,
		nx,
		ny,
		100,
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

	tempCategory := ""
	for _, item := range items.Item {
		// 각 카테고리 별로 가장 먼저 들어온 데이터 저장
		if tempCategory == item.Category {
			continue
		} else {
			tempCategory = item.Category
		}

		temp := entity.WhetherSrtEntity{}
		temp.Key = item.Category
		if item.Category == "VEC" {
			temp.Value = entity.WhetherVecString(item.FcstValue)
		} else {
			temp.Value = item.FcstValue
		}
		whetherRes = append(whetherRes, temp)
	}

	return whetherRes, nil
}

// func: 기상청 기상특보통보문 api 조회
// @param
func (s *ServiceWhether) GetWhetherWrnMsg() (entity.WhetherWrnMsgList, error) {

	now := time.Now()
	startDate := now.AddDate(0, 0, -6).Format("20060102")
	endDate := now.Format("20060102")

	url := fmt.Sprintf("http://apis.data.go.kr/1360000/WthrWrnInfoService/getWthrWrnMsg?serviceKey=%s&pageNo=%s&numOfRows=%s&dataType=%s&fromTmFc=%s&toTmFc=%s&stnId=%s",
		s.ApiKey.WhetherApiKey, // 데이터 포털 API 키
		"1",                    // 페이지 번호
		"20",                   // 한 페이지 결과 수
		"JSON",                 // 응답 자료 형식
		startDate,              // 발표시각 from
		endDate,                //endDate,                // 발표시각 to
		"108") // stnId. 전국(108), 서울(109), 부산(159), 대구(143), 광주(156), 전주(146), 대전(133), 청주(131), 강릉(105), 제주(184)

	// api call
	body, err := api.CallGetAPI(url)
	if err != nil {
		return nil, fmt.Errorf("call GetWhetherWrnMsgList API error: %v", err)
	}

	// api response item struct
	type WhetherMsg struct {
		Msg string `json:"t6"`
	}

	type whether struct {
		Response struct {
			Header struct {
				ResultCode string `json:"resultCode"`
				ResultMsg  string `json:"resultMsg"`
			} `json:"header"`
			Body struct {
				Items struct {
					Item []WhetherMsg `json:"item"`
				} `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	// response parse
	var res whether
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, fmt.Errorf("WhetherWrnMsg api JSON parse err: %v", err)
	}

	if res.Response.Header.ResultCode != "00" {
		return nil, fmt.Errorf("WhetherWrnMsg api response err : %s", res.Response.Header.ResultMsg)
	}

	// whether api response -> go api response convert
	items := res.Response.Body.Items
	msg := items.Item[0].Msg

	// 데이터가 문자열로 길게 와서 각 특보별, 지역별 구분하여 데이터 반환
	// 데이터 예시: "o 건조주의보 : 경기도(안산, 시흥, 김포, 평택, 화성 제외), 강원도(강릉평지, 동해평지, 태백, 삼척평지, 속초평지, 고성평지, 양양평지, 영월, 정선평지, 원주, 강원남부산지), 충청북도, 전라남도(곡성, 구례, 고흥, 보성, 여수, 광양, 순천, 장흥), 전북자치도(무주, 남원), 경상북도, 경상남도(통영, 고성 제외), 서울, 대전, 광주, 대구, 부산, 울산"
	list := entity.WhetherWrnMsgList{}
	warningMsgs := strings.Split(msg, "o")
	warningList := []string{"강풍주의보", "강풍경보", "호우주의보", "호우경보", "대설주의보", "대설경보", "태풍주의보", "태풍경보", "황사주의보", "황사경보", "폭염주의보", "폭염경보", "한파주의보", "한파경보"} // , "건조주의보"

	for _, warningMsg := range warningMsgs {
		response := entity.WhetherWrnMsg{}
		warning, areaStr, resultBool := strings.Cut(warningMsg, ":")

		warning = strings.TrimSpace(warning)

		// 제공할 특보인 경우에만 진행 - 제공 특보(warningList): 강풍, 호우, 대설, 태풍, 황사, 폭염, 한파
		flag := false
		for _, warningStr := range warningList {
			if warningStr == warning {
				flag = true
				break
			}
		}
		// 발효된 특보가 제공할 특보가 아닌 경우
		if flag == false {
			continue
		}

		// 특보 정보 입력
		response.Warning = warning
		var areaList []string

		// 특보 지역 구분(시도 단위로)
		if resultBool {
			areaStr = strings.ReplaceAll(areaStr, "도, ", "도), ")
			for _, area := range strings.Split(areaStr, "), ") {
				// 값이 없으면 넘기기
				if area == "" {
					continue
				}
				addStr := ""
				if strings.Contains(area, "(") {
					addStr = ")"
				}
				areaList = append(areaList, strings.TrimSpace(area)+addStr)
			}
		} else {
			continue
		}

		list = append(list, response)
	}

	return list, nil
}
