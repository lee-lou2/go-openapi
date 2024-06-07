package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

// GetDB 데이터베이스 인스턴스
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		serverEnv := GetEnv("SERVER_ENV")

		if serverEnv == "prod" || serverEnv == "stag" || serverEnv == "dev" {
			var (
				host     = GetEnv("DB_HOST")
				port     = GetEnv("DB_PORT")
				user     = GetEnv("DB_USER")
				password = GetEnv("DB_PASSWORD")
				dbname   = GetEnv("DB_NAME")
			)

			dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				log.Fatalf("failed to connect database: %v", err)
			}

			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("failed to get generic database object: %v", err)
			}

			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(time.Hour)

			dbInstance = db
		} else {
			dbFileName := "sqlite.db"
			if IsTesting() {
				dbFileName = "test.sqlite.db"
			}
			db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				log.Fatalf("failed to connect database: %v", err)
			}

			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("failed to get generic database object: %v", err)
			}

			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(time.Hour)

			dbInstance = db
		}
	})

	return dbInstance
}
