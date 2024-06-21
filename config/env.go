package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
)

func init() {
	// 환경 변수 불러오기
	_ = godotenv.Load()
	// 파라미터 스토어 조회
	_ = GetAWSParams()
}

// GetEnv 환경 변수 조회
func GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv 환경 변수 설정
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// GetAWSParams AWS 파라미터 스토어에서 환경 변수 조회
func GetAWSParams() error {
	if GetEnv("AWS_ACCESS_KEY_ID") == "" && GetEnv("AWS_SECRET_ACCESS_KEY") == "" {
		return nil
	}
	ctx := context.TODO()
	awsRegion := GetEnv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return err
	}
	svc := ssm.NewFromConfig(cfg)
	serverEnv := GetEnv("SERVER_ENV")
	path := fmt.Sprintf("/go-openapi/%s/", serverEnv)
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}
	result, err := svc.GetParametersByPath(ctx, input)
	if err != nil {
		return err
	}
	for _, param := range result.Parameters {
		key := strings.TrimPrefix(*param.Name, path)
		value := *param.Value
		os.Setenv(key, value)
	}
	return nil
}
