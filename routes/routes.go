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
	return router
}
