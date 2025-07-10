package repositories

import (
	"log"

	"gorm.io/gorm"

	"github.com/nabilulilalbab/toko-klontong/config"
	"github.com/nabilulilalbab/toko-klontong/models"
)

func GetAllProduk() []models.Produk {
	var produk []models.Produk
	err := config.DB.Find(&produk)
	if err != nil {
		log.Println("REPOSITORI: gagal get all produk")
	}
	return produk
}

// start new

type ProdukRepository interface {
	FindByID(id uint) (models.Produk, error)
	FindAll() ([]models.Produk, error)
	Save(produk models.Produk) (models.Produk, error)
	Update(produk models.Produk) (models.Produk, error)
	Delete(id uint) error
	SearchByName(name string) ([]models.Produk, error)
}

type produkRepositoryImpl struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepositoryImpl{db: db}
}

func (r *produkRepositoryImpl) FindAll() ([]models.Produk, error) {
	var produks []models.Produk
	err := r.db.Find(&produks).Error
	if err != nil {
		return nil, err
	}
	return produks, nil
}

func (r *produkRepositoryImpl) Save(produk models.Produk) (models.Produk, error) {
	if err := r.db.Create(&produk).Error; err != nil {
		return models.Produk{}, err
	}
	return produk, nil
}

func (r *produkRepositoryImpl) FindByID(id uint) (models.Produk, error) {
	var produk models.Produk
	err := r.db.First(&produk, id).Error
	return produk, err
}

func (r *produkRepositoryImpl) Update(produk models.Produk) (models.Produk, error) {
	err := r.db.Save(&produk).Error
	return produk, err
}

func (r *produkRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Produk{}, id).Error
}

func (r *produkRepositoryImpl) SearchByName(name string) ([]models.Produk, error) {
	var produks []models.Produk
	// Menggunakan kueri LIKE dengan wildcard '%' untuk mencari bagian dari nama
	err := r.db.Where("nama_produk LIKE ?", "%"+name+"%").Find(&produks).Error
	return produks, err
}
