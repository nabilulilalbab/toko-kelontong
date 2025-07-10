package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/nabilulilalbab/toko-klontong/models"
)

type TransaksiRepository interface {
	Create(transaksi models.Transaksi) (models.Transaksi, error)
	FindAll() ([]models.Transaksi, error)
	FindByID(id uint) (models.Transaksi, error)
}

type transaksiRepositoryImpl struct {
	db *gorm.DB
}

func NewTransaksiRepository(db *gorm.DB) TransaksiRepository {
	return &transaksiRepositoryImpl{db: db}
}

func (r *transaksiRepositoryImpl) Create(transaksi models.Transaksi) (models.Transaksi, error) {
	// GORM Transaction: Memastikan semua operasi di dalamnya berhasil atau di-rollback (dibatalkan) semua.
	// Ini penting untuk menjaga konsistensi data.
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan data transaksi utama. GORM akan otomatis menyimpan DetailTransaksis
		//    yang ada di dalamnya dan mengisi TransaksiID.
		if err := tx.Create(&transaksi).Error; err != nil {
			return err // Rollback jika gagal
		}

		// 2. Loop melalui setiap item detail untuk mengurangi stok produk.
		for _, item := range transaksi.DetailTransaksis {
			var produk models.Produk

			// Ambil data produk yang akan dikurangi stoknya.
			if err := tx.First(&produk, item.ProdukID).Error; err != nil {
				return err // Rollback jika produk tidak ditemukan
			}

			// Cek ketersediaan stok.
			if produk.Stok < item.Jumlah {
				return errors.New("stok tidak mencukupi untuk produk: " + produk.NamaProduk) // Rollback
			}

			// Kurangi stok dan simpan kembali.
			produk.Stok -= item.Jumlah
			if err := tx.Save(&produk).Error; err != nil {
				return err // Rollback jika gagal update stok
			}
		}

		// Jika tidak ada error sama sekali, kembalikan nil untuk COMMIT transaksi.
		return nil
	})
	if err != nil {
		return models.Transaksi{}, err
	}

	return transaksi, nil
}

// FindAll untuk mengambil semua data transaksi
func (r *transaksiRepositoryImpl) FindAll() ([]models.Transaksi, error) {
	var transaksis []models.Transaksi
	// Urutkan berdasarkan yang terbaru
	err := r.db.Order("created_at desc").Find(&transaksis).Error
	return transaksis, err
}

// FindByID untuk mengambil satu transaksi beserta detailnya
func (r *transaksiRepositoryImpl) FindByID(id uint) (models.Transaksi, error) {
	var transaksi models.Transaksi
	// GORM Preload untuk otomatis mengambil data relasi (DetailTransaksis dan Produk-nya)
	err := r.db.Preload("DetailTransaksis.Produk").First(&transaksi, id).Error
	return transaksi, err
}
