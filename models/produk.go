package models

import "gorm.io/gorm"

type Produk struct {
	ID         uint   `gorm:"primarykey`
	NamaProduk string `gorm:not null`
	Harga      uint   `gorm:not null`
	Stok       uint   `gorm:not null;default:0`
	gorm.Model        // Ini akan otomatis menambahkan field CreatedAt, UpdatedAt, DeletedAt
}
