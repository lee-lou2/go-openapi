# Golang Open API 🚀

## 프로젝트 소개
이 프로젝트는 Golang을 이용하여 Open API를 개발하는 개인 프로젝트입니다. 아래와 같은 과제를 수행하면서 기술적인 도전과 문제 해결을 목표로 하고 있습니다.

## 과제
- **인증**: OAuth2.0 기반으로 개발자 등록(Authorization Code), API(Client Credentials) 구현
- **권한 및 제한**: scopes를 통한 권한 부여 및 분당 Call Count 관리(등급제 운영)
- **개발자 페이지**: Client Key 발급 및 사용량 확인을 위한 페이지 제공
- **관리자 페이지**: 계정별 권한 설정 및 전체 통계 확인 페이지 구축
- **보안**: 개인 정보 암호화, XSS, SQL 인젝션 등 보안 위협 요소 차단
- **기타**: AWS 파라미터 스토어를 이용한 환경 변수 관리

## 어떤 서비스를 제공할 것인가?
- **고민중입니다 🧐**

## 사용 기술 스택
- **언어**: Golang
- **인증**: OAuth2.0(Authorization Code, Client Credentials)
- **데이터베이스**: Sqlite3(Local), MySQL(QA, Prod)
- **클라우드**: AWS

## 폴더 구조(작성중)
```
/go-openapi
├── api
│   ├── handler
│   │   ├── auth.go
│   │   └── user.go
│   ├── middleware
│   │   ├── auth.go
│   │   └── logger.go
│   ├── router
│   │   └── v1.go
│   └── server.go
├── db
│   ├── migrations
│   └── models
├── internal
│   ├── auth
│   ├── logger
│   ├── server
│   └── user
├── pkg
│   ├── auth
│   ├── db
│   ├── logger
│   └── user
├── configs
│   └── config.go
├── scripts
├── docs
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## 실행 방법
### 사전 준비
- Go 설치 (>= 1.21)
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

3. Docker 컨테이너 실행
   ```sh
   docker-compose up --build
   ```

4. 개발 서버 실행
   ```sh
   go run main.go
   ```
