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
