package services

import (
	"project/models"
	"project/repositories"

	"gorm.io/gorm"
)

func AddOne(id int) (int, error) {
	return (id + 1), nil
}

type TestService struct {
	DB *gorm.DB
}

func NewTestService(db *gorm.DB) *TestService {
	return &TestService{DB: db}
}

func (s *TestService) CreateTest(test *models.Test) error {
	return repositories.CreateTest(s.DB, test)
}

func (s *TestService) GetTests() (*[]models.Test, error) {
	return repositories.GetTests(s.DB)
}

func (s *TestService) DeleteTest(id uint) error {
	return repositories.DeleteTest(s.DB, id)
}

func (s *TestService) GetTestByID(id uint) (*models.Test, error) {
	return repositories.GetTestByID(s.DB, id)
}

func (s *TestService) RawQuery() (*[]models.Test, error) {
	return repositories.RawQuery(s.DB)
}
