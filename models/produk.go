package models

import "gorm.io/gorm"

type Produk struct {
	ID         uint   `gorm:"primaryKey"`
	NamaProduk string `gorm:"not null"`
	Harga      uint   `gorm:"not null"`
	Stok       uint   `gorm:"not null;default:0"`
	gorm.Model        // Ini akan otomatis menambahkan field CreatedAt, UpdatedAt, DeletedAt
}
