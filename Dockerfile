# 빌드 단계
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# 실행 단계
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY . .
EXPOSE 80
COPY --from=builder /app/main /app/main
CMD ["/app/main"]