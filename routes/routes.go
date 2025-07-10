package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/controllers"
)

func SetupRouter(produkController controllers.ProdukController, transaksiController controllers.TransaksiController) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Selamat Datang di Toko Klontong!")
	})
	router.GET("/produk", produkController.Index)
	router.GET("/produk/tambah", produkController.Form)
	router.POST("/produk/tambah", produkController.Store)
	router.GET("/produk/edit/:id", produkController.Edit)
	router.POST("/produk/update/:id", produkController.Update)
	router.POST("/produk/delete/:id", produkController.Delete)
	router.GET("/produk/export", produkController.Export)
	router.GET("/kasir", transaksiController.ShowKasirPage)
	router.POST("/kasir/checkout", transaksiController.ProcessCheckout)
	router.GET("/histori", transaksiController.ShowHistoryPage)
	router.GET("/reports/histori", transaksiController.ExportPDF)
	router.GET("/histori/:id", transaksiController.ShowHistoryDetailPage)
	router.GET("/api/produk/search", produkController.SearchAPI)
	router.POST("/api/generate-qris", transaksiController.GenerateQRIS)
	router.GET("/pembayaran/:id", transaksiController.ShowPaymentPage)
	router.POST("/pembayaran/konfirmasi/:id", transaksiController.ConfirmPayment)
	return router
}
