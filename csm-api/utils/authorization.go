package utils

import (
	"github.com/guregu/null"
	"golang.org/x/exp/slices"
	"strings"
)

// func: 권한 슬라이스에서 role이 있는지 확인 후 있으면 ture값 반환, 없으면 false값 반환
// @param
// - list: 권한 담긴 슬라이스
// - roles: 역할 문자열 (|로 나열된 것)
func AuthorizationListCheck(list []string, roles null.String) bool {
	check := false

	for _, role := range strings.Split(roles.String, "|") {
		role = strings.TrimSpace(role)
		if slices.Contains(list, role) {
			check = true
		}
	}

	return check
}
