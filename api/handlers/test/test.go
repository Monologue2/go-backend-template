package handlers

import (
	"net/http"
	"strconv"

	"project/api/services"
	"project/models"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// TestHandler 구조체
type TestHandler struct {
	TestService *services.TestService
}

// TestHandler 생성자
func NewTestHandler(service *services.TestService) *TestHandler {
	return &TestHandler{TestService: service}
}

// 사용자 생성
func (h *TestHandler) CreateTest(c *gin.Context) {
	var test models.Test

	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TestService.CreateTest(&test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Test created successfully", "test": test})
}

// 사용자 목록 조회
func (h *TestHandler) GetTests(c *gin.Context) {
	tests, err := h.TestService.GetTests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tests"})
		return
	}
	c.JSON(http.StatusOK, tests)
}

// 특정 사용자 조회
func (h *TestHandler) GetTestByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	test, err := h.TestService.GetTestByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}
	c.JSON(http.StatusOK, test)
}

// 사용자 삭제
func (h *TestHandler) DeleteTest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	if err := h.TestService.DeleteTest(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test deleted successfully"})
}

func (h *TestHandler) AddOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	addedId, err := services.AddOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}
	c.JSON(http.StatusOK, addedId)
}

func (h *TestHandler) RawQuery(c *gin.Context) {
	test, err := h.TestService.RawQuery()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed RawQuery Services."})
	}
	c.JSON(http.StatusOK, test)
}
