package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/nabilulilalbab/toko-klontong/models"
)

type ProdukRepositoryMock struct {
	mock.Mock
}

func (m *ProdukRepositoryMock) Save(produk models.Produk) (models.Produk, error) {
	args := m.Called(produk)
	return args.Get(0).(models.Produk), args.Error(1)
}

func (m *ProdukRepositoryMock) FindAll() ([]models.Produk, error) {
	args := m.Called()
	return args.Get(0).([]models.Produk), args.Error(1)
}
