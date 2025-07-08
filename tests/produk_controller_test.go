package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/nabilulilalbab/toko-klontong/controllers"
	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/tests/mocks"
	"github.com/nabilulilalbab/toko-klontong/utils"
)

func TestProdukController_Index_Success(t *testing.T) {
	mockService := new(mocks.ProdukServiceMock)

	// Setup mock
	expectedProduks := []models.Produk{
		{NamaProduk: "Susu Bendera", Harga: 1500, Stok: 20},
	}
	mockService.On("FindAll").Return(expectedProduks, nil)

	templates := utils.ParseTemplates()
	produkcontroller := controllers.NewProdukController(mockService, templates)

	request := httptest.NewRequest("GET", "/produk", nil)
	recorder := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/produk", produkcontroller.Index)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Susu Bendera")

	// Verifikasi bahwa mock method benar-benar dipanggil
	mockService.AssertExpectations(t)
}
