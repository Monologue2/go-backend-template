# Go Backend Template
[Web server gin](https://go.dev/doc/tutorial/web-service-gin)
> 핸들러 + 서비스 + 저장소(Repository) 패턴: Golang에서 MVC 패턴을 단순화한 구조<br>
> 개인적으로 프로젝트 양산 겸 사용을 위해 작성하였습니다.

- DB_TYPE
- GIN_MODE

## Directory Structure
| 계층         | 역할 |
|-------------|----------------------------------|
| **`api/services/`** | 비즈니스 로직 실행, 여러 Repository 조합 가능 |
| **`api/handlers/`** | Context 해석, 요청 검증, Service 실행 후 클라이언트 응답 |
| **`api/routes/`** | Handler와 HTTP 엔드포인트 연결 |
| **`middlewares/`** | 요청 로깅, 인증/인가 등 공통 로직 관리 |
| **`models/`** | Struct 정의, 테이블 마이그레이션 |
| **`repositories/`** | DB 연결 관리, ORM 또는 SQL 실행 |
| **`config/`** | YAML 기반 설정 관리 (환경별 Config 지원) |
| **`main.go`** | 전체 애플리케이션 초기화 및 실행 |

## *gin.Context
| 기능 | 설명 | 예제 |
|------|------|------|
| **요청 데이터 읽기** | URL 파라미터, 쿼리 스트링, 바디 데이터 등을 가져올 수 있음 | `c.Query("name")`, `c.Param("id")`, `c.PostForm("email")` |
| **응답 보내기** | JSON, XML, HTML 등의 형태로 클라이언트에 응답 | `c.JSON(200, gin.H{"message": "ok"})` |
| **미들웨어 간 데이터 공유** | `c.Set()`으로 저장하고 `c.Get()`으로 가져옴 | `c.Set("user", userData)`, `c.Get("user")` |
| **상태 코드 설정** | HTTP 응답 상태 코드 변경 | `c.Status(400)`, `c.AbortWithStatus(403)` |
| **파일 업로드 처리** | `c.FormFile()`을 이용해 파일을 업로드 | `file, _ := c.FormFile("file")` |

## Flow

### models/에서 테이블 정의
```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name  string `gorm:"type:varchar(100);not null"`
    Email string `gorm:"uniqueIndex;not null"`
}
```

### models/migrate.go에서 DB 마이그레이션 적용
```go
package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
    db.AutoMigrate(&User{})
}
```

### repositories/에서 Repository 계층 작성
```go
package repositories

import (
    "gorm.io/gorm"
    "project/models"
)

func CreateUser(db *gorm.DB, user *models.User) error {
    return db.Create(user).Error
}

func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
    var user models.User
    result := db.First(&user, id)
    return &user, result.Error
}
```

### services/에서 비즈니스 로직 및 DB 연결 관리
```go
package services

import (
    "gorm.io/gorm"
    "project/models"
    "project/repositories"
)

type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) error {
    return repositories.CreateUser(s.DB, user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
    return repositories.GetUserByID(s.DB, id)
}
```

### handlers/에서 Service를 호출하여 HTTP 요청 처리
```go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "project/models"
    "project/services"
)

type UserHandler struct {
    Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
    return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Service.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}
```

### routes/에서 API 엔드포인트 정의
```go
package routes

import (
    "github.com/gin-gonic/gin"
    "project/handlers"
)

func SetupUserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
    v1 := r.Group("/api/v1/users")
    {
        v1.POST("/", userHandler.CreateUser)
        v1.GET("/:id", userHandler.GetUserByID)
    }
}
```

### main.go에서 모든 모듈을 초기화하고 서버 실행
```go
package main

import (
    "github.com/gin-gonic/gin"
    "project/repositories"
    "project/services"
    "project/handlers"
    "project/routes"
    "project/models"
)

func main() {
    r := gin.Default()

    // DB 초기화
    db := repositories.DB
    models.Migrate(db)

    // 서비스 & 핸들러 초기화
    userService := services.NewUserService(db)
    userHandler := handlers.NewUserHandler(userService)

    // 라우트 설정
    routes.SetupUserRoutes(r, userHandler)

    // 서버 실행
    r.Run(":8080")
}
```

## 개선점
[] JWT 인증/인가 추가 → middlewares/에서 AuthMiddleware 추가
[] Swagger 연동 → docs/ 폴더 추가 후 OpenAPI 문서화
[] gRPC 지원 → grpc/ 디렉터리 추가 후 protobuf 정의
[] Pub/Sub 메시징 (Kafka, RabbitMQ) → services/에서 메시지 큐 연동
[] 테스트 코드 추가 (*_test.go) → repository_test.go, service_test.go, handler_test.go 등