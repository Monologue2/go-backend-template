package repositories

import (
	"project/models"

	"gorm.io/gorm"
)

func CreateTest(db *gorm.DB, test *models.Test) error {
	return db.Create(test).Error
}

func GetTests(db *gorm.DB) (*[]models.Test, error) {
	var tests []models.Test
	result := db.Find(&tests)
	return &tests, result.Error
}

func DeleteTest(db *gorm.DB, id uint) error {
	return db.Delete(id).Error
}

func GetTestByID(db *gorm.DB, id uint) (*models.Test, error) {
	var test models.Test
	result := db.First(&test, id)
	return &test, result.Error
}

func RawQuery(db *gorm.DB) (*[]models.Test, error) {
	var testValue []models.Test
	result := db.Raw("SELECT id, test FROM tests").Scan(&testValue) // 이것도 repo로 옮길 것
	return &testValue, result.Error
}
