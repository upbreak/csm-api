# Go 빌드 환경 설정
FROM golang:1.23.4 AS builder

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 복사 및 종속성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사 및 빌드
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app

# 실행 환경 설정
FROM base-dependencies

# Go 애플리케이션 복사
COPY --from=builder /app/app /app/

# 작업 디렉토리 설정
WORKDIR /app

# 실행 명령어
CMD ["./app"]
