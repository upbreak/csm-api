//go:build prod

package config

import (
	"log"
	"os"
)

func init() {
	log.Println("start go:build prod")

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

	err = os.Setenv("UPLOAD_PATH", "/var/csm/uploads")
	if err != nil {
		log.Fatalf("Failed to set UPLOAD_PATH: %v", err)
	}

	err = os.Setenv("LOG_PATH", "/var/csm/logs")
	if err != nil {
		log.Fatalf("Failed to set LOG_PATH: %v", err)
	}

	err = os.Setenv("ERR_LOG_PATH", "/var/csm/logs/error")
	if err != nil {
		log.Fatalf("Failed to set LOG_PATH: %v", err)
	}

	err = os.Setenv("EXCEL_PATH", "/var/csm/resources/excel")
	if err != nil {
		log.Fatalf("Failed to set EXCEL_PATH: %v", err)
	}

	err = os.Setenv("CONSOLE_LOG_PATH", "/var/log/csm")
	if err != nil {
		log.Fatalf("Failed to set EXCEL_PATH: %v", err)
	}
}
