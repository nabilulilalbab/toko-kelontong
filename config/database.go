package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nabilulilalbab/toko-klontong/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("toko.db"), &gorm.Config{})
	if err != nil {
		panic("Gagal terhubung ke database!")
	}
	err = db.AutoMigrate(&models.Produk{})
	if err != nil {
		panic("Gagal melakukan migrasi database!")
	}
	DB = db
}
