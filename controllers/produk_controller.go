package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/services"
)

// controllers new start

type ProdukController interface {
	Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type produkcotrollerimpl struct {
	produkService services.ProdukService
	templates     *template.Template
}

func NewProdukController(s services.ProdukService, t *template.Template) ProdukController {
	return &produkcotrollerimpl{produkService: s, templates: t}
}

func (c *produkcotrollerimpl) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Controller: Memulai proses menampilkan daftar produk.")
	produks, err := c.produkService.FindAll()
	if err != nil {
		log.Printf("Controller: Error saat mengambil data produk: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := map[string]any{
		"Title":    "Daftar Produk",
		"Products": produks,
	}
	err = c.templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Controller: Error saat merender template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Println("Controller: Berhasil menampilkan daftar produk.")
}
