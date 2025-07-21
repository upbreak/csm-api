# Go 빌드 환경 설정
FROM golang:1.23.4 AS builder

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 복사 및 종속성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사 및 빌드
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags prod -o app

# 실행 환경 설정
FROM base-dependencies

# 빌드 시점에 사용할 고정 UID/GID를 인자로 받습니다 (기본값 101)
ARG CIMAGE_UID=101
ARG CIMAGE_GID=101

# non-root 그룹 및 사용자 생성 (고정된 UID/GID 사용)
RUN addgroup --system --gid ${CIMAGE_GID} cimage \
 && adduser --system --uid ${CIMAGE_UID} --gid ${CIMAGE_GID} --disabled-password --gecos "" cimage

# Go 애플리케이션 복사
COPY --from=builder /app/app /app/

# 작업 디렉토리 설정
WORKDIR /app

# 애플리케이션 소유권을 고정 UID/GID의 cimage 사용자로 변경
RUN chown ${CIMAGE_UID}:${CIMAGE_GID} /app/app

# non-root 사용자로 전환
USER cimage

# 실행 명령어
CMD ["./app"]
