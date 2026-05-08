package database

import (
	"reserach/pkg/model"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&model.Product{}).Count(&count)

	if count > 0 {
		return
	}

	user := model.User{Username: "benchmark_user", Email: "test@lab.com"}
	db.Create(&user)

	product := model.Product{Name: "Benchmark Product", Price: 999.99, Stock: 50000}

	product.ID = 1
	db.Create(&product)

	println("Seed Data wurde angelegt: Product ID 1 mit Stock 50000")
}
