//go:build dev

package config

import (
	"log"
	"os"
)

func init() {
	// 개발 환경에 필요한 기본값 설정
	err := os.Setenv("ENV", "development")
	if err != nil {
		log.Fatalf("Failed to set ENV: %v", err)
	}

	err = os.Setenv("PORT", "8082")
	if err != nil {
		log.Fatalf("Failed to set PORT: %v", err)
	}
}
