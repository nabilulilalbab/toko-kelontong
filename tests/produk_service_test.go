package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/services"
	"github.com/nabilulilalbab/toko-klontong/tests/mocks"
)

func TestProdukService_FindAll_Success(t *testing.T) {
	// 1. Setup
	mockRepo := new(mocks.ProdukRepositoryMock)
	// Data produk yang kita harapkan akan dikembalikan oleh repository
	expectedProduks := []models.Produk{
		{ID: 1, NamaProduk: "Indomie Goreng", Harga: 3000, Stok: 100},
	}
	// 2. "Program" mock: Jika method FindAll() dipanggil, kembalikan 'expectedProduks' dan tanpa error.
	mockRepo.On("FindAll").Return(expectedProduks, nil)
	// Buat instance service dengan mock repository
	produkService := services.NewProdukService(mockRepo)
	// 3. Eksekusi
	actualProduks, err := produkService.FindAll()
	// 4. Assert (Validasi)
	assert.Nil(t, err)                              // Pastikan tidak ada error
	assert.Equal(t, expectedProduks, actualProduks) // Pastikan data yang kembali sesuai harapan
	mockRepo.AssertExpectations(t)                  // Pastikan method yang di-mock benar-benar dipanggil
}

func TestProdukService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.ProdukRepositoryMock)
	inputProduk := models.Produk{
		NamaProduk: "Kopi Kapal Api",
		Harga:      2000,
		Stok:       150,
	}
	mockRepo.On("Save", inputProduk).Return(inputProduk, nil)
	produkService := services.NewProdukService(mockRepo)
	createdProduk, err := produkService.Create(inputProduk)
	assert.Nil(t, err)
	assert.Equal(t, "Kopi Kapal Api", createdProduk.NamaProduk)
	mockRepo.AssertExpectations(t)
}
