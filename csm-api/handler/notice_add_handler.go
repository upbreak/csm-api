package handler

import (
	"csm-api/auth"
	"csm-api/service"
)

// 제목, 내용,, request에서 받은 정보 service로 넘겨서, db에 저장하기
// service

type NoticeAddHandler struct {
	Service service.NoticeService
	Jwt     *auth.JWTUtils
}
