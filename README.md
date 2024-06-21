# Golang Open API 🚀

## 프로젝트 소개
Golang 자체 net/http 를 이용하여 Open API 를 개발하는 개인 프로젝트

## 서비스 및 구성
- **사용자**: 사용자 생성 및 관련 기능
- **인증**: OAuth2.0 기반으로 개발자 등록(Bearer JWT), API(Client Credentials) 구현
- **권한 및 제한**: scopes, user level 을 통한 권한 부여 및 분당 Call Count 관리(등급제 운영)
- **사용자 페이지**: 회원 가입, 로그인, 패스워드 재설정 및 클라이언트 관리
- **보안**: 개인 정보(email) 암호화, 보안 위협 요소 차단
- **미들웨어**: static, cors, recover, request_id 등
- **기타**: AWS 파라미터 스토어를 이용한 환경 변수 관리

## 사용 기술 스택
- **언어**: Golang
- **인증**: Bearer(JWT), OAuth2.0(Client Credentials)
- **데이터베이스**: Sqlite3(Local), Postgresql(DEV, QA, Prod)

## 테스트
- httptest 를 이용한 핸들러 검증(test 패키지)
- 각 기능별 유닛테스트

## 폴더 구조
```
/go-openapi
├── api
│   ├── handler
│   │   ├── auth
│   │   │   ├── client.go
│   │   │   ├── login.go
│   │   │   └── token.go
│   │   ├── client
│   │   │   └── me.go
│   │   └── user
│   │       ├── password.go
│   │       ├── user.go
│   │       └── verify.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── chain.go
│   │   ├── cors.go
│   │   ├── limit.go
│   │   ├── logger.go
│   │   ├── permission.go
│   │   ├── recover.go
│   │   └── request_id.go
│   ├── parser
│   │   └── json.go
│   ├── render
│   │   └── json.go
│   ├── router
│   │   ├── base.go
│   │   ├── template.go
│   │   └── v1.go
│   └── validation
│       ├── oauth2.go
│       ├── user.go
│       └── verify.go
├── cmd
│   └── api
│       └── server.go
├── configs
│   ├── cache.go
│   ├── db.go
│   ├── env.go
│   └── test.go
├── internal
│   ├── auth
│   │   ├── client.go
│   │   ├── login.go
│   │   └── token.go
│   └── user
│       ├── password.go
│       ├── user.go
│       └── verify.go
├── model
│   ├── client
│   │   ├── client.go
│   │   ├── const.go
│   │   └── enum.go
│   └── user
│       └── user.go
├── pkg
│   ├── auth
│   │   ├── client.go
│   │   ├── token_test.go
│   │   └── token.go
│   ├── notify
│   │   └── email.go
│   ├── user
│   │   └── verify.go
│   └── utils
│       ├── base64_test.go
│       ├── base64.go
│       ├── encrypt_test.go
│       ├── encrypt.go
│       ├── hash_test.go
│       ├── hash.go
│       ├── password_test.go
│       ├── password.go
│       ├── string_test.go
│       └── string.go
├── test
│   ├── .env
│   ├── handler_auth_test.go
│   ├── handler_user_test.go
│   └── server_test.go
├── views
│   ├── accounts
│   │   ├── login
│   │   │   └── index.html
│   │   ├── password
│   │   │   └── index.html
│   │   └── signup
│   │       └── index.html
│   ├── user
│   │   ├── password
│   │   │   └── index.html
│   │   └── verify
│   │       └── index.html
│   ├── error
│   │   └── index.html
│   └── index.html
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── deploy.sh
├── Dockerfile
├── main.go
└── README.md
```

## 실행 방법
### 사전 준비
- Go 설치 (>= 1.22)
- Docker 설치
- 환경 변수 설정 (AWS 파라미터 스토어 사용)

### 실행 단계
1. 저장소 클론
   ```sh
   git clone https://github.com/lee-lou2/go-openapi
   cd go-openapi
   ```

2. 종속성 설치
   ```sh
   go mod tidy
   ```

3. 개발 서버 실행
   ```sh
   go run .
   ```

### 환경 변수
- 로컬 실행 시 `.env` 파일 필요
```sh
SERVER_ENV=
# 파라미터 스토어 사용시
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=ap-northeast-2
# 파라미터 스토어 사용하지 않을시
SERVER_HOST=
SERVER_PORT=
EMAIL_SMTP_HOST=
EMAIL_SMTP_PORT=
EMAIL_USERNAME=
EMAIL_PASSWORD=
JWT_SECRET=
COOKIE_ENCRYPT_KEY=
SHA256_SALT=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
AES256_KEY=
```