package routes

import (
	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/controllers"
)

func SetupRoutes(router *httprouter.Router) {
	router.GET("/produk", controllers.ProdukIndex)
	// router.GET("/produk/create", controllers.ProdukCreateForm)
	// router.POST("/product", controllers.ProdukStore)
	// router.GET("/produk/edit/:id", controllers.ProdukEditForm)
	// router.POST("/produk/update/:id", controllers.ProdukUpdate)
	// router.POST("/produk/delete/:id", controllers.ProdukDelete)
}
