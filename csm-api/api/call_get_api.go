package api

import (
	"csm-api/utils"
	"io"
	"net/http"
)

func CallGetAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", utils.CustomMessageErrorf("GET 요청 실패", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", utils.CustomMessageErrorf("응답 읽기 실패", err)
	}

	return string(body), nil
}
