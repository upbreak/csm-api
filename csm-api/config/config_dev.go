//go:build dev

package config

import (
	"log"
	"os"
)

func init() {
	log.Println("start go:build dev")

	// 개발 환경에 필요한 기본값 설정
	err := os.Setenv("ENV", "development")
	if err != nil {
		log.Fatalf("Failed to set ENV: %v", err)
	}

	err = os.Setenv("ROLE", "web")
	if err != nil {
		log.Fatalf("Failed to set ROLE: %v", err)
	}

	err = os.Setenv("DOMAIN", "127.0.0.1")
	if err != nil {
		log.Fatalf("Failed to set DOMAIN: %v", err)
	}

	err = os.Setenv("PORT", "8082")
	if err != nil {
		log.Fatalf("Failed to set PORT: %v", err)
	}

	err = os.Setenv("UPLOAD_PATH", "uploads")
	if err != nil {
		log.Fatalf("Failed to set UPLOAD_PATH: %v", err)
	}
}
