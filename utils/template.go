package utils

import (
	"html/template"
	"log"
)

func ParseTemplates() *template.Template {
	log.Println("Parsing templates...")

	// 1. Parsing semua file di folder 'templates' (untuk layout.html)
	tmplt, err := template.ParseGlob("../templates/*.html")
	if err != nil {
		log.Fatalf("Gagal mem-parsing template layout: %v", err)
	}

	// 2. Tambahkan parsing untuk file di dalam sub-folder 'produk'
	// Perhatikan kita menggunakan 'tmplt.ParseGlob' (bukan template.ParseGlob)
	// untuk menambahkan ke template set yang sudah ada.
	_, err = tmplt.ParseGlob("../templates/produk/*.html")
	if err != nil {
		log.Fatalf("Gagal mem-parsing template produk: %v", err)
	}

	log.Println("Parsing templates selesai.")
	return tmplt
}

