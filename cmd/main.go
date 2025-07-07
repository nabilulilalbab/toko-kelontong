package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/config"
	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/repositories"
	"github.com/nabilulilalbab/toko-klontong/routes"
)

func main() {
	// database
	config.ConnectDatabase()
	fmt.Println("koneksi database Berhasil")
	// instance repository
	produkRepo := repositories.NewProdukRepository(config.DB)
	seedProducts(produkRepo)
	// router
	router := httprouter.New()
	routes.SetupRoutes(router)
	fmt.Println("Server berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

// Fungsi untuk seeding data produk
func seedProducts(repo repositories.ProdukRepository) {
	// Cek apakah produk sudah ada, jika belum, tambahkan.
	produks, _ := repo.FindAll()
	if len(produks) == 0 {
		fmt.Println("Melakukan seeding data produk...")
		repo.Save(models.Produk{NamaProduk: "Indomie Goreng", Harga: 3000, Stok: 100})
		repo.Save(models.Produk{NamaProduk: "Teh Botol Sosro", Harga: 5000, Stok: 50})
		repo.Save(models.Produk{NamaProduk: "Silverqueen 65g", Harga: 12500, Stok: 30})
		fmt.Println("Seeding selesai.")
	}
}
