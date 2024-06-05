# 빌드 단계
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main

# 실행 단계
FROM alpine:latest

RUN apk --no-cache add ca-certificates

EXPOSE 80

COPY --from=build /app/main /main

CMD ["/main"]
