package service

import (
	"context"
	"csm-api/api"
	"csm-api/config"
	"csm-api/entity"
	"csm-api/store"
	"csm-api/txutil"
	"csm-api/utils"
	"encoding/json"
	"fmt"
	"github.com/guregu/null"
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
type ServiceWeather struct {
	ApiKey       *config.ApiConfig
	SafeDB       store.Queryer
	SafeTDB      store.Beginner
	Store        store.WeatherStore
	SitePosStore store.SitePosStore
}

// func: 기상청 초단기예보 api 조회
// @param
// - date string: 현재날짜(mmdd), time string: 현재시간(hhmm), nx int: 위도변환값, ny int: 경도변환값
func (s *ServiceWeather) GetWeatherSrtNcst(date string, time string, nx int, ny int) (entity.WeatherSrtEntityRes, error) {
	// 초단기실황 url
	url := fmt.Sprintf("http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getUltraSrtFcst?dataType=JSON&ServiceKey=%s&base_date=%s&base_time=%s&nx=%d&ny=%d&numOfRows=%d",
		s.ApiKey.DataGoApiKey,
		date,
		time,
		nx,
		ny,
		100,
	)

	// api call
	body, err := api.CallGetAPI(url)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	// api response item struct
	type Weather struct {
		Response struct {
			Body struct {
				Items entity.WeatherSrtItems `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	// response parse
	var res Weather
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, utils.CustomErrorf(err)
	}

	// Weather api response -> go api response convert
	items := res.Response.Body.Items
	weatherRes := entity.WeatherSrtEntityRes{}

	tempCategory := ""
	for _, item := range items.Item {
		// 각 카테고리 별로 가장 먼저 들어온 데이터 저장
		if tempCategory == item.Category {
			continue
		} else {
			tempCategory = item.Category
		}

		temp := entity.WeatherSrtEntity{}
		temp.Key = item.Category
		if item.Category == "VEC" {
			temp.Value = entity.WeatherVecString(item.FcstValue)
		} else {
			temp.Value = item.FcstValue
		}
		weatherRes = append(weatherRes, temp)
	}

	return weatherRes, nil
}

// func: 기상청 기상특보통보문 api 조회
// @param
func (s *ServiceWeather) GetWeatherWrnMsg() (entity.WeatherWrnMsgList, error) {

	now := time.Now()
	startDate := now.AddDate(0, 0, -6).Format("20060102")
	endDate := now.Format("20060102")

	url := fmt.Sprintf("http://apis.data.go.kr/1360000/WthrWrnInfoService/getWthrWrnMsg?serviceKey=%s&pageNo=%s&numOfRows=%s&dataType=%s&fromTmFc=%s&toTmFc=%s&stnId=%s",
		s.ApiKey.DataGoApiKey, // 데이터 포털 API 키
		"1",                   // 페이지 번호
		"20",                  // 한 페이지 결과 수
		"JSON",                // 응답 자료 형식
		startDate,             // 발표시각 from
		endDate,               //endDate,                // 발표시각 to
		"108")                 // stnId. 전국(108), 서울(109), 부산(159), 대구(143), 광주(156), 전주(146), 대전(133), 청주(131), 강릉(105), 제주(184)

	// api call
	body, err := api.CallGetAPI(url)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	// api response item struct
	type WeatherMsg struct {
		Msg string `json:"t6"`
	}

	type Weather struct {
		Response struct {
			Header struct {
				ResultCode string `json:"resultCode"`
				ResultMsg  string `json:"resultMsg"`
			} `json:"header"`
			Body struct {
				Items struct {
					Item []WeatherMsg `json:"item"`
				} `json:"items"`
			} `json:"body"`
		} `json:"response"`
	}

	// response parse
	var res Weather
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, utils.CustomErrorf(err)
	}

	if res.Response.Header.ResultCode != "00" {
		return nil, utils.CustomErrorf(fmt.Errorf("WeatherWrnMsg api response err : %s", res.Response.Header.ResultMsg))
	}

	// Weather api response -> go api response convert
	items := res.Response.Body.Items
	msg := items.Item[0].Msg

	// 데이터가 문자열로 길게 와서 각 특보별, 지역별 구분하여 데이터 반환
	// 데이터 예시: "o 건조주의보 : 경기도(안산, 시흥, 김포, 평택, 화성 제외), 강원도(강릉평지, 동해평지, 태백, 삼척평지, 속초평지, 고성평지, 양양평지, 영월, 정선평지, 원주, 강원남부산지), 충청북도, 전라남도(곡성, 구례, 고흥, 보성, 여수, 광양, 순천, 장흥), 전북자치도(무주, 남원), 경상북도, 경상남도(통영, 고성 제외), 서울, 대전, 광주, 대구, 부산, 울산"
	list := entity.WeatherWrnMsgList{}
	warningMsgs := strings.Split(msg, "o")
	warningList := []string{"강풍주의보", "강풍경보", "호우주의보", "호우경보", "대설주의보", "대설경보", "태풍주의보", "태풍경보", "황사주의보", "황사경보", "폭염주의보", "폭염경보", "한파주의보", "한파경보"} // , "건조주의보"

	for _, warningMsg := range warningMsgs {
		response := entity.WeatherWrnMsg{}
		warning, areaStr, _ := strings.Cut(warningMsg, ":")

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

		var areaList []string

		// 특보 지역 구분(시도 단위로) 완도, 진도, 초도 등의 지명인 경우 이상하게 나와서 보류.
		//if resultBool {
		//
		//	areaStr = strings.ReplaceAll(areaStr, "도, ", "도), ")
		//	for _, area := range strings.Split(areaStr, "), ") {
		//		// 값이 없으면 넘기기
		//		if area == "" {
		//			continue
		//		}
		//		addStr := ""
		//		if strings.Contains(area, "(") {
		//			addStr = ")"
		//		}
		//		areaList = append(areaList, strings.TrimSpace(area)+addStr)
		//	}
		//
		//} else {
		//	continue
		//}
		areaList = append(areaList, areaStr)

		// 특보 정보 입력
		if len(areaList) > 0 {
			response.Warning = warning
			response.Area = areaList
			list = append(list, response)
		}
	}
	return list, nil
}

// func: 저장된 날씨 리스트 조회
// params
// - sno: 현장 PK
// - targetDate: 조회할 날짜
func (s *ServiceWeather) GetWeatherList(ctx context.Context, sno int64, targetDate time.Time) (*entity.Weathers, error) {

	weathers, err := s.Store.GetWeatherList(ctx, s.SafeDB, sno, targetDate)
	if err != nil {
		return nil, utils.CustomErrorf(err)
	}

	return weathers, nil
}

// func: IRIS_SITE_POS에 저장된 현장 날씨 저장(스케줄러)
// params:
// -
func (s *ServiceWeather) SaveWeather(ctx context.Context) (err error) {
	tx, cleanup, err := txutil.BeginTxWithCleanMode(ctx, s.SafeTDB, false)
	if err != nil {
		return utils.CustomErrorf(err)
	}

	defer func() {
		txutil.DeferTx(tx, &err)
		cleanup()
	}()

	// IRIS_SITE_POS에 등록된 값들 불러오기
	list, err := s.SitePosStore.GetSitePosList(ctx, s.SafeDB)

	if err != nil {
		err = utils.CustomErrorf(err)
	}

	for _, site := range list {

		// 장소를 바꿀 수 없는 경우
		if site.Latitude.Valid == false || site.Longitude.Valid == false {
			continue
		}

		now := time.Now()
		baseDate := now.Format("20060102")
		baseTime := now.Add(time.Minute * -30).Format("1504") // 기상청에서 30분 단위로 발표하기 때문에 30분 전의 데이터 요청
		nx, ny := utils.LatLonToXY(site.Latitude.Float64, site.Longitude.Float64)

		// 초단기예보 조회
		res, weatherErr := s.GetWeatherSrtNcst(baseDate, baseTime, nx, ny)
		if weatherErr != nil {
			err = utils.CustomMessageErrorf("not json format", weatherErr)
		}

		// weather 형태로 변경
		weather, convertErr := s.convertWeather(res)

		if convertErr != nil {
			err = utils.CustomErrorf(convertErr)
		}

		// 값이 없는 경우 저장하지 않음
		if !weather.Lgt.Valid || !weather.Pty.Valid || !weather.Sky.Valid || !weather.Rn1.Valid || !weather.T1h.Valid || !weather.Wsd.Valid || !weather.Vec.Valid {
			continue
		}

		weather.Sno = site.Sno
		weather.RecogTime = null.TimeFrom(now)

		// weather 저장
		if err = s.Store.SaveWeather(ctx, tx, *weather); err != nil {
			err = utils.CustomErrorf(err)
		}
	}

	return
}

// func: 초단기예보 key, value로 구성된 데이터를 Weather객체로 변경
func (s *ServiceWeather) convertWeather(res entity.WeatherSrtEntityRes) (*entity.Weather, error) {
	weather := &entity.Weather{}
	for _, weatherSrt := range res {
		if weatherSrt.Key == "LGT" {
			weather.Lgt = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "PTY" {
			weather.Pty = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "RN1" {
			weather.Rn1 = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "SKY" {
			weather.Sky = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "T1H" {
			weather.T1h = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "REH" {
			weather.Reh = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "UUU" {
			weather.Uuu = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "VVV" {
			weather.Vvv = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "VEC" {
			weather.Vec = utils.ParseNullString(weatherSrt.Value)
		} else if weatherSrt.Key == "WSD" {
			weather.Wsd = utils.ParseNullString(weatherSrt.Value)
		}
	}

	return weather, nil
}
