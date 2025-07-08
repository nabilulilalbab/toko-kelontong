package services

import (
	"log"

	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/repositories"
)

// new services start

type ProdukService interface {
	FindAll() ([]models.Produk, error)
	Create(produk models.Produk) (models.Produk, error)
}

type produkServiceImpl struct {
	repo repositories.ProdukRepository
}

func (s *produkServiceImpl) FindAll() ([]models.Produk, error) {
	log.Println("Service: Memulai proses FindAll produk.")
	produks, err := s.repo.FindAll()
	if err != nil {
		log.Printf("Service: Error saat memanggil repo.FindAll(): %v\n", err)
		return nil, err
	}
	log.Println("Service: Berhasil mendapatkan data produk.")
	return produks, nil
}

func NewProdukService(repo repositories.ProdukRepository) ProdukService {
	return &produkServiceImpl{repo: repo}
}

func (s *produkServiceImpl) Create(produk models.Produk) (models.Produk, error) {
	log.Println("Service: Memulai proses create produk.")
	newProduk, err := s.repo.Save(produk)
	if err != nil {
		log.Printf("Service: Error saat menyimpan produk baru: %v\n", err)
		return models.Produk{}, err
	}
	return newProduk, nil
}
