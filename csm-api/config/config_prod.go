//go:build prod

package config

import (
	"log"
	"os"
)

func init() {
	log.Println("start go:build dev")

	// 운영 환경에 필요한 기본값 설정
	err := os.Setenv("ENV", "production")
	if err != nil {
		log.Fatalf("Failed to set ENV: %v", err)
	}

	//err = os.Setenv("DOMAIN", "csm.htenc.co.kr")
	err = os.Setenv("DOMAIN", "0.0.0.0")
	if err != nil {
		log.Fatalf("Failed to set DOMAIN: %v", err)
	}

	err = os.Setenv("PORT", "8080")
	if err != nil {
		log.Fatalf("Failed to set PORT: %v", err)
	}
}
