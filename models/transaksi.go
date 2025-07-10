package models

import "gorm.io/gorm"

type Transaksi struct {
	ID               uint   `gorm:"primaryKey"`
	TotalHarga       uint   `gorm:"not null"`
	MetodePembayaran string `gorm:"not null"`
	Status           string `gorm:"not null;default:'pending'"`
	gorm.Model

	DetailTransaksis []DetailTransaksi `gorm:"foreignKey:TransaksiID"`
}
type DetailTransaksi struct {
	ID          uint `gorm:"primaryKey"`
	TransaksiID uint `gorm:"not null"` // Foreign key ke tabel Transaksi
	ProdukID    uint `gorm:"not null"` // Foreign key ke tabel Produk
	Jumlah      uint `gorm:"not null"`
	Subtotal    uint `gorm:"not null"`

	// Relasi: Menghubungkan ke struct Produk untuk mempermudah pengambilan data nama produk
	Produk Produk `gorm:"foreignKey:ProdukID"`
}
