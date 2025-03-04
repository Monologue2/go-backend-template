package repositories

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 구조체
type Config struct {
	App struct {
		Mode         string `yaml:"mode"`
		GormLogLevel string `yaml:"gorm_log_level"`
	} `yaml:"app"`

	Postgres struct {
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		User            string `yaml:"user"`
		Password        string `yaml:"password"`
		DBName          string `yaml:"dbname"`
		SSLMode         string `yaml:"sslmode"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
		ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
		ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
	} `yaml:"postgres"`
}

// 전역 설정 변수
var AppConfig Config
var DB *gorm.DB

func init() {
	if err := LoadConfig(); err != nil {
		log.Fatalf("❌ Config load failed: %v", err)
	}

	var err error
	DB, err = InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
}

func LoadConfig() error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	filename := fmt.Sprintf("config/config-%s.yaml", env)

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		return fmt.Errorf("failed to decode config file: %v", err)
	}

	os.Setenv("GIN_MODE", AppConfig.App.Mode)

	log.Printf("✅ Loaded config: %s\n", filename)
	return nil
}

func getGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent // 기본값 (최소 로그)
	}
}

// DB 초기화
func InitDB() (*gorm.DB, error) {
	switch os.Getenv("DB_TYPE") {
	case "postgres":
		return connectPostgres()
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s", os.Getenv("DB_TYPE"))
	}
}

// PostgreSQL 연결
func connectPostgres() (*gorm.DB, error) {
	cfg := AppConfig.Postgres

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getGormLogLevel(AppConfig.App.GormLogLevel)), // 로그 레벨 설정
	})
	if err != nil {
		log.Fatalf("❌ Database connection is not initialized\nPlease check your configuration file or ENV variable.\n")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Database connection is not initialized\nPlease check your postgreSQL driver.\n")
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("✅ Connected to PostgreSQL(GORM)")
	return db, nil
}
