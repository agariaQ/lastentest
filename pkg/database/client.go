package database

import (
	"fmt"
	"log"
	"reserach/pkg/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", "localhost", 5432, "postgres", "postgres", "postgres")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Konnte nicht zur Datenbank verbinden:", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{})
	if err != nil {
		log.Fatal("Konnte Tabellen nicht automatisch migrieren:", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Datenbankverbindung steht und Tabellen sind migriert.")
	return db
}
