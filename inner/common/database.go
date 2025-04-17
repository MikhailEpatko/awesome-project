package common

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ConnectDb подключиться к базе данных и вернуть подключение
func ConnectDb() *sqlx.DB {
	cfg := GetConfig(".env")
	return ConnectDbWithCfg(cfg)
}

// ConnectDbWithCfg подключиться к базе данных с использованием конфигурации и вернуть подключение
func ConnectDbWithCfg(cfg Config) *sqlx.DB {
	return sqlx.MustConnect(cfg.DbDriverName, cfg.Dsn)
}
