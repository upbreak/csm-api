//go:build prod

package config

import (
	"log"
	"os"
)

func init() {
	// 운영 환경에 필요한 기본값 설정
	err := os.Setenv("ENV", "production")
	if err != nil {
		log.Fatalf("Failed to set ENV: %v", err)
	}

	err = os.Setenv("PORT", "8080")
	if err != nil {
		log.Fatalf("Failed to set PORT: %v", err)
	}
}
