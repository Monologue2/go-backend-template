package models

import "gorm.io/gorm"

// AutoMigrate 실행 (DB에 테이블 자동 생성)
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Test{})
}
