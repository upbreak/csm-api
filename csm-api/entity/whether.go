package entity

import "strconv"

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-03-05
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct: 초단기 예보 api 응답 구조
type WhetherSrtItem struct {
	BaseDate  string `json:"baseDate"`
	BaseTime  string `json:"baseTime"`
	Category  string `json:"category"`
	ObsrValue string `json:"obsrValue"`
}
type WhetherSrtItems struct {
	Item []WhetherSrtItem `json:"item"`
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
//   - VEC	풍향	° (16방위) (WhetherVecString() 참고)
//   - WSD	풍속	m/s (미터/초)
type WhetherSrtEntity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type WhetherSrtEntityRes []WhetherSrtEntity

// func: 풍향 변환 숫자 -> 한글
// @param
// - value: 풍향(방위)
func WhetherVecString(value string) string {
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
