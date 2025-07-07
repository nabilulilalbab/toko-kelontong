package models

import "gorm.io/gorm"

type Produk struct {
	ID         uint    `gorm:"primarykey`
	NamaProduk string  `gorm:not null`
	Harga      float64 `gorm:not null`
	Stok       int     `gorm:not null;default:0`
	gorm.Model         // Ini akan otomatis menambahkan field CreatedAt, UpdatedAt, DeletedAt
}
