package db

import (
	"fmt"
	"log"

	"github.com/Aller101/calc-back-go/internal/config"
	"github.com/Aller101/calc-back-go/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(c config.Config) (*gorm.DB, error) {
	dsn := DSN(c)

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not con to db: %v", err)
	}
	if err := db.AutoMigrate(&service.Calculation{}); err != nil {
		log.Fatalf("Migrate: %v", err)
	}
	return db, nil
}

func DSN(c config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Dbname,
		c.Port,
		c.Sslmode,
	)
}
