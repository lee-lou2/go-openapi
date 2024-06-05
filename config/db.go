package config

import (
	"log"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbFileName = "sqlite.db"

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

// GetDB 데이터베이스 인스턴스
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
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
	})

	return dbInstance
}
