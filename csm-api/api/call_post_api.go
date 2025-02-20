package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CallPostAPI(url string, payload interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("JSON 변환 실패: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("POST 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("응답 읽기 실패: %w", err)
	}

	return string(body), nil
}
