package common

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

// ConnectDb подключиться к базе данных и вернуть подключение
func ConnectDb() *sqlx.DB {
	cfg := GetConfig(".env")
	return ConnectDbWithCfg(cfg)
}

// ConnectDbWithCfg подключиться к базе данных с использованием конфигурации и вернуть подключение
func ConnectDbWithCfg(cfg Config) *sqlx.DB {
	var db = sqlx.MustConnect(cfg.DbDriverName, cfg.Dsn)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}
