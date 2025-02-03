package auth

import (
	"context"
	"csm-api/clock"
	"csm-api/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JWTUtils struct {
	Cfg   *config.JwtConfig
	Clock clock.Clocker
}

type JWTRole string

const (
	Admin JWTRole = "admin"
	User  JWTRole = "user"
)

// claims 정의
type JWTClaims struct {
	UserId string
	Role   JWTRole
	Token  string
}

// context에 저장하는 claims value
type UserId struct{}
type Role struct{}

// JWTUtils 구조체 생성
func JwtNew(c clock.Clocker) (*JWTUtils, error) {
	jwt := &JWTUtils{}

	jwtConfig, err := config.GetJwtConfig()
	if err != nil {
		return nil, err
	}

	jwt.Cfg = jwtConfig
	jwt.Clock = c

	return jwt, nil
}

// 토큰 생성
func (j *JWTUtils) GenerateToken(jwtClaims *JWTClaims) (string, error) {
	// 비밀 키
	secretKey := []byte(j.Cfg.SecretKey)

	// JWT 클레임 설정
	claims := jwt.MapClaims{
		"userId":     jwtClaims.UserId,
		"exp":        j.Clock.Now().Add(1 * time.Minute).Unix(), // 만료 시간: 1분
		"role":       jwtClaims.Role,
		"refreshKey": j.Cfg.RefreshKey,
	}

	// 토큰 생성
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 비밀 키로 서명
	tokenString, err := parseToken.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("jwtUtils.go/GenerateToken() err: %w", err)
	}

	jwtClaims.Token = tokenString

	return tokenString, nil
}

// 토큰 재발급
func (j *JWTUtils) RefreshToken(refreshToken string) (*JWTClaims, error) {
	// 비밀 키
	secretKey := []byte(j.Cfg.SecretKey)

	// 토큰 파싱 및 검증
	parseToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// 서명 방법 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwtUtile.go/RefreshToken-invalid refreshToken: %v", err)
	}

	// 클레임 RefreshKey 확인
	//claims, ok := parseToken.Claims.(jwt.MapClaims)
	claims, _ := parseToken.Claims.(jwt.MapClaims)
	if claims["refreshKey"] != j.Cfg.RefreshKey {
		return nil, fmt.Errorf("jwtUtile.go/RefreshToken-refreshToken is differ: %v", err)
	}
	//if !ok || !parseToken.Valid {
	//	return nil, fmt.Errorf("jwtUtile.go/RefreshToken-invalid token")
	//}

	// JWT 클레임 설정
	jwtClaims := jwt.MapClaims{
		"userId":     claims["userId"].(string),
		"exp":        j.Clock.Now().Add(1 * time.Minute).Unix(), // 만료 시간: 1분
		"role":       claims["role"].(string),
		"refreshKey": j.Cfg.RefreshKey,
	}

	// 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// 비밀 키로 서명
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return &JWTClaims{}, fmt.Errorf("jwtUtils.go/GenerateToken() err: %w", err)
	}

	jwtTokenClaims := &JWTClaims{
		UserId: claims["userId"].(string),
		Role:   JWTRole(claims["role"].(string)),
		Token:  tokenString,
	}

	return jwtTokenClaims, nil
}

// 토큰 유효성 검사
func (j *JWTUtils) ValidateJWT(r *http.Request) (*JWTClaims, error) {
	// 비밀 키
	secretKey := []byte(j.Cfg.SecretKey)
	//fmt.Printf("validateJWT secretkey: %v \n", secretKey)

	// 쿠키 읽기
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, fmt.Errorf("jwtUtils.go/validateJWT() err: %v", err)
	}
	tokenString := cookie.Value

	// 토큰 파싱 및 검증
	parseToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 서명 방법 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwtUtile.go/invalid token: %v", err)
	}

	// 클레임 확인
	claims, ok := parseToken.Claims.(jwt.MapClaims)
	if !ok || !parseToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	jwtClaims := &JWTClaims{
		UserId: claims["userId"].(string),
		Role:   JWTRole(claims["role"].(string)),
	}

	return jwtClaims, nil
}

// 사용자 토큰 인증시 필요한 api 호출시 token의 claims에 있는 값을 context에 저장
func (j *JWTUtils) FillContext(r *http.Request) (*http.Request, error) {
	// 토큰 검사 및 claims 추출
	claims, err := j.ValidateJWT(r)
	if err != nil {
		return nil, err
	}

	// claims 데이터 context에 저장
	ctx := SetContext(r.Context(), UserId{}, claims.UserId)
	ctx = SetContext(ctx, Role{}, string(claims.Role))

	httpRequestClone := r.Clone(ctx)

	return httpRequestClone, nil
}

func SetContext(ctx context.Context, key struct{}, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContext(ctx context.Context, key interface{}) (string, bool) {
	value, ok := ctx.Value(key).(string)
	return value, ok
}
