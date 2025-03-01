package models

import "gorm.io/gorm"

// (Mapping to Test Table)
type Test struct {
	ID   uint   `gorm:"primaryKey"` // 기본 키
	Test string `gorm:"type:varchar(100);not null"`
}

// AutoMigrate 실행 (DB에 테이블 자동 생성)
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Test{})
}
