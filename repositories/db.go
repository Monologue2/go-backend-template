package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL 드라이버
	"gopkg.in/yaml.v3"
)

// Config 구조체
type Config struct {
	App struct {
		Mode string `yaml:"mode"`
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

func init() {
	if err := LoadConfig(); err != nil {
		log.Fatalf("❌ Config load failed: %v", err)
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

// DB 초기화
func InitDB() (*sql.DB, error) {
	switch os.Getenv("DB_TYPE") {
	case "postgres":
		return connectPostgres()
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s", os.Getenv("DB_TYPE"))
	}
}

// PostgreSQL 연결
func connectPostgres() (*sql.DB, error) {
	host := AppConfig.Postgres.Host
	port := AppConfig.Postgres.Port
	user := AppConfig.Postgres.User
	password := AppConfig.Postgres.Password
	dbname := AppConfig.Postgres.DBName
	sslmode := AppConfig.Postgres.SSLMode

	psqlDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlDsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}
