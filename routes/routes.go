package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/controllers"
)

func SetupRouter(produkController controllers.ProdukController) *httprouter.Router {
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
	router.GET("/produk/export/excel", produkController.ExportExcel)
	return router
}
