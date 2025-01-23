module csm-api

go 1.23.4

require (
	github.com/caarlos0/env v3.5.0+incompatible // indirect:: struct에 태그를 사용하여 환경변수를 가져와 사용할 수 있게 해준다
	github.com/go-chi/chi/v5 v5.2.0 // indirect:: net/http 패키지 타입 정의를 따르며 라우팅 기능을 제공
	github.com/godror/godror v0.46.0 // indirect:: Oracle 데이터베이스와의 연결을 지원하는 Go 라이브러리
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect:: JSON Web Token(JWT)을 생성하고, 파싱하며, 검증하는 데 사용되는 라이브러리
	github.com/jmoiron/sqlx v1.4.0 // indirect:: database/sql 패키지를 확장한 라이브러리로, SQL 쿼리 실행 및 결과 매핑을 보다 직관적이고 간편하게 할 수 있도록 도와준다.
	github.com/rs/cors v1.11.1 // indirect:: CORS(Cross-Origin Resource Sharing) 를 처리할 수 있도록 도와주는 라이브러리
	golang.org/x/sync v0.10.0 // indirect:: main()에서 사용할 run함수를 구현할 때 이용하는 준표준 패키지
)

require (
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/godror/knownpb v0.2.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	google.golang.org/protobuf v1.36.2 // indirect
)
