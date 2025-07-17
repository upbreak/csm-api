package api

import (
	"bytes"
	"csm-api/utils"
	"encoding/json"
	"io"
	"net/http"
)

func CallPostAPI(url string, payload interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", utils.CustomMessageErrorf("JSON 변환 실패", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", utils.CustomMessageErrorf("POST 요청 실패", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", utils.CustomMessageErrorf("응답 읽기 실패", err)
	}

	return string(body), nil
}
