package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/nabilulilalbab/toko-klontong/models"
)

type ProdukServiceMock struct {
	mock.Mock
}

func (m *ProdukServiceMock) FindAll() ([]models.Produk, error) {
	args := m.Called()
	return args.Get(0).([]models.Produk), args.Error(1)
}

func (m *ProdukServiceMock) Create(produk models.Produk) (models.Produk, error) {
	args := m.Called(produk)
	return args.Get(0).(models.Produk), args.Error(1)
}

func (m *ProdukServiceMock) Update(id uint, produk models.Produk) (models.Produk, error) {
	args := m.Called(id, produk)
	return args.Get(0).(models.Produk), args.Error(1)
}

func (m *ProdukServiceMock) GetByID(id uint) (models.Produk, error) {
	args := m.Called(id)
	return args.Get(0).(models.Produk), args.Error(1)
}

func (m *ProdukServiceMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProdukServiceMock) Search(query string) ([]models.Produk, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Produk), args.Error(1)
}
