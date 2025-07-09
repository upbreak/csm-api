package entity

import (
	"github.com/guregu/null"
	"strconv"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-03-05
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 초단기 예보 api 응답 구조
type WeatherSrtItem struct {
	FcstDate  string `json:"fcstDate"`
	FcstTime  string `json:"fcstTime"`
	Category  string `json:"category"`
	FcstValue string `json:"fcstValue"`
}
type WeatherSrtItems struct {
	Item []WeatherSrtItem `json:"item"`
}

// struct: 날씨 예보 최종 응답
// @Key-Value:
//   - T1H	기온	°C (섭씨)
//   - RN1	1시간 강수량	mm (밀리미터)
//   - REH	습도	% (퍼센트)
//   - PTY	강수 형태	코드 값
//     0	없음 (비 안 옴)
//     1	비
//     2	비 또는 눈 (진눈깨비)
//     3	눈
//     4	소나기
//     5	빗방울
//     6	빗방울 또는 눈날림
//     7	눈날림
//   - SKY	하늘 상태	코드 값
//     1	맑음
//     3	구름 많음
//     4	흐림
//   - VEC	풍향	° (16방위) (WeatherVecString() 참고)
//   - WSD	풍속	m/s (미터/초)
type WeatherSrtEntity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type WeatherSrtEntityRes []WeatherSrtEntity

type WeatherSrt struct {
	Weather WeatherSrtEntityRes `json:"weather"`
	Sno     int64               `json:"sno"`
}
type WeatherSrtRes []WeatherSrt

type Weather struct {
	Sno       null.Int    `json:"sno" db:"SNO"`               // 현장번호
	Lgt       null.String `json:"lgt" db:"LGT"`               // 낙뢰
	Pty       null.String `json:"pty" db:"PTY"`               // 강수형태
	Rn1       null.String `json:"rn1" db:"RN1"`               // 1시간 강수량
	Sky       null.String `json:"sky" db:"SKY"`               // 하늘상태
	T1h       null.String `json:"t1h" db:"T1H"`               // 기온
	Reh       null.String `json:"reh" db:"REH"`               // 습도
	Uuu       null.String `json:"uuu" db:"UUU"`               // 풍속 - 동서 성분
	Vvv       null.String `json:"vvv" db:"VVV"`               // 풍속 - 남북성분
	Vec       null.String `json:"vec" db:"VEC"`               // 풍향(숫자)
	VecKo     null.String `json:"vec_ko" db:"VEC_KO"`         // 풍향(한글)
	Wsd       null.String `json:"wsd" db:"WSD"`               // 풍속
	RecogTime null.Time   `json:"recog_time" db:"RECOG_TIME"` // 날씨 측정 시각
}

type Weathers []*Weather

// func: 풍향 변환 숫자 -> 한글
// @param
// - value: 풍향(방위)
func WeatherVecString(value string) string {
	vec, _ := strconv.ParseFloat(value, 64)
	if vec > 337.5 || vec <= 22.5 {
		return "북(N)"
	} else if vec > 22.5 && vec <= 67.5 {
		return "북동(NE)"
	} else if vec > 67.5 && vec <= 112.5 {
		return "동(E)"
	} else if vec > 112.5 && vec <= 157.5 {
		return "남동(SE)"
	} else if vec > 157.5 && vec <= 202.5 {
		return "남(S)"
	} else if vec > 202.5 && vec <= 247.5 {
		return "남서(SW)"
	} else if vec > 247.5 && vec <= 292.5 {
		return "서(W)"
	} else if vec > 292.5 && vec <= 337.5 {
		return "북서(NW)"
	} else {
		return "없음"
	}
}

type WeatherWrnMsg struct {
	Warning string   `json:"warning"`
	Area    []string `json:"area"`
}

type WeatherWrnMsgList []WeatherWrnMsg
