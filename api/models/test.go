package models

// (Mapping to Test Table)
type Test struct {
	ID   uint   `gorm:"primaryKey"` // 기본 키
	Test string `gorm:"type:varchar(100);not null"`
}
