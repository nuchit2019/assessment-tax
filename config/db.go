package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Config struct {
	Db            *sql.DB
	Port          string
	Admin         string
	AdminPassword string
}

func New() (*Config, error) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return &Config{
		Db:            db,
		Port:          os.Getenv("PORT"),
		Admin:         os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}, nil
}

func (c *Config) Close() error {
	if c.Db != nil {
		return c.Db.Close()
	}
	return nil
}
