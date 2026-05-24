package database

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sea-cucumber-trace/backend/internal/config"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.DBDriver)) {
	case "", "sqlite", "sqlite3":
		return gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	case "mysql":
		return gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q", cfg.DBDriver)
	}
}
