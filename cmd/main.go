package main

import (
	"fmt"
	"net/http"

	"github.com/nabilulilalbab/toko-klontong/config"
	"github.com/nabilulilalbab/toko-klontong/controllers"
	"github.com/nabilulilalbab/toko-klontong/repositories"
	"github.com/nabilulilalbab/toko-klontong/routes"
	"github.com/nabilulilalbab/toko-klontong/services"
	"github.com/nabilulilalbab/toko-klontong/utils"
)

func main() {
	config.ConnectDatabase()
	// Inisialisasi semua lapisan (Dependency Injection)
	templateCache := utils.ParseTemplates()
	produkRepo := repositories.NewProdukRepository(config.DB)
	produkService := services.NewProdukService(produkRepo)
	produkController := controllers.NewProdukController(produkService, templateCache)
	// Inisialisasi dependensi untuk Transaksi
	transaksiRepo := repositories.NewTransaksiRepository(config.DB)
	transaksiService := services.NewTransaksiService(transaksiRepo, produkRepo)
	transaksiController := controllers.NewTransaksiController(transaksiService, produkService, templateCache)
	router := routes.SetupRouter(produkController, transaksiController)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	fmt.Println("Server berjalan di http://localhost:8080")
	fmt.Println("Akses daftar produk di http://localhost:8080/produk")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
