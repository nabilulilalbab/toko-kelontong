package services

import (
	"log"

	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/repositories"
)

// CreateTransaksiRequest adalah struct untuk menampung data input dari keranjang belanja
type CreateTransaksiRequest struct {
	Items []struct {
		ProdukID uint `json:"produk_id"`
		Jumlah   uint `json:"jumlah"`
	} `json:"items"`
}

type TransaksiService interface {
	Create(request CreateTransaksiRequest) (models.Transaksi, error)
	GetAll() ([]models.Transaksi, error)
	GetByID(id uint) (models.Transaksi, error)
}

type transaksiServiceImpl struct {
	transaksiRepo repositories.TransaksiRepository
	produkRepo    repositories.ProdukRepository
}

// NewTransaksiService membutuhkan dua repository sebagai dependensi
func NewTransaksiService(transaksiRepo repositories.TransaksiRepository, produkRepo repositories.ProdukRepository) TransaksiService {
	return &transaksiServiceImpl{
		transaksiRepo: transaksiRepo,
		produkRepo:    produkRepo,
	}
}

func (s *transaksiServiceImpl) Create(request CreateTransaksiRequest) (models.Transaksi, error) {
	log.Println("Service: Memulai proses create transaksi.")

	var totalHarga uint
	var detailTransaksis []models.DetailTransaksi

	// Loop melalui setiap item di keranjang untuk validasi dan perhitungan
	for _, item := range request.Items {
		// Ambil data produk dari DB untuk mendapatkan harga yang valid
		produk, err := s.produkRepo.FindByID(item.ProdukID)
		if err != nil {
			return models.Transaksi{}, err
		}

		subtotal := produk.Harga * item.Jumlah
		totalHarga += subtotal

		detailTransaksis = append(detailTransaksis, models.DetailTransaksi{
			ProdukID: item.ProdukID,
			Jumlah:   item.Jumlah,
			Subtotal: subtotal,
		})
	}

	// Siapkan data transaksi utama
	transaksi := models.Transaksi{
		TotalHarga:       totalHarga,
		DetailTransaksis: detailTransaksis,
	}

	// Panggil repository untuk menyimpan semuanya dalam satu DB Transaction
	newTransaksi, err := s.transaksiRepo.Create(transaksi)
	if err != nil {
		return models.Transaksi{}, err
	}

	log.Println("Service: Berhasil membuat transaksi baru.")
	return newTransaksi, nil
}

func (s *transaksiServiceImpl) GetAll() ([]models.Transaksi, error) {
	log.Println("Service: Mengambil semua data transaksi.")
	return s.transaksiRepo.FindAll()
}

func (s *transaksiServiceImpl) GetByID(id uint) (models.Transaksi, error) {
	log.Printf("Service: Mencari transaksi dengan ID: %d", id)
	return s.transaksiRepo.FindByID(id)
}
