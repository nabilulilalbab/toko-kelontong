package services

import (
	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/repositories"
)

func GetProdukList() []models.Produk {
	return repositories.GetAllProduk()
}
