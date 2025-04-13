package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

type Config struct {
	dbDriverName string
	dsn          string
}

func Connect() *sqlx.DB {
	cfg := getConfig(".env")
	return ConnectWithCfg(cfg)
}

func ConnectWithCfg(cfg Config) *sqlx.DB {
	return sqlx.MustConnect(cfg.dbDriverName, cfg.dsn)
}

func getConfig(envFile string) Config {
	var err = godotenv.Load(envFile)
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err.Error()))
	}
	return Config{
		dbDriverName: os.Getenv("DB_DRIVER_NAME"),
		dsn:          os.Getenv("DB_DSN"),
	}
}
