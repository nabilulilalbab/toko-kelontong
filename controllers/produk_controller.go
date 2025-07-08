package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/services"
)

// controllers new start

type ProdukController interface {
	Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Form(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type produkControllerImpl struct {
	produkService services.ProdukService
	templates     *template.Template
}

func NewProdukController(s services.ProdukService, t *template.Template) ProdukController {
	return &produkControllerImpl{produkService: s, templates: t}
}

func (c *produkControllerImpl) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	err = c.templates.ExecuteTemplate(w, "list.html", data)
	if err != nil {
		log.Printf("Controller: Error saat merender template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Println("Controller: Berhasil menampilkan daftar produk.")
}

func (c *produkControllerImpl) Form(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := map[string]any{
		"Title":      "Tambah Produk Baru",
		"FormAction": "/produk/tambah/",
	}
	if err := c.templates.ExecuteTemplate(w, "form.html", data); err != nil {
		log.Printf("Controller: Error saat merender template form: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *produkControllerImpl) Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Controller: Memulai proses store produk.")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal mem-parsing form", http.StatusBadRequest)
		return
	}
	namaProduk := r.PostForm.Get("nama_produk")
	harga, _ := strconv.Atoi(r.PostForm.Get("harga"))
	stok, _ := strconv.Atoi(r.PostForm.Get("stok"))
	produk := models.Produk{
		NamaProduk: namaProduk,
		Harga:      uint(harga),
		Stok:       uint(stok),
	}
	_, err := c.produkService.Create(produk)
	if err != nil {
		http.Error(w, "Gagal menyimpan produk", http.StatusInternalServerError)
		return
	}
	log.Println("Controller: Berhasil store produk, redirecting...")
	http.Redirect(w, r, "/produk", http.StatusSeeOther)
}

func (c *produkControllerImpl) Edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	produk, err := c.produkService.GetByID(uint(id))
	if err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}
	data := map[string]any{
		"Title":      "Edit Produk",
		"Produk":     produk,
		"FormAction": "/produk/update/" + strconv.Itoa(id),
	}
	c.templates.ExecuteTemplate(w, "form.html", data)
}

func (c *produkControllerImpl) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	r.ParseForm()

	produk := models.Produk{
		NamaProduk: r.PostForm.Get("nama_produk"),
		Harga:      uint(mustAtoi(r.PostForm.Get("harga"))),
		Stok:       uint(mustAtoi(r.PostForm.Get("stok"))),
	}

	c.produkService.Update(uint(id), produk)
	http.Redirect(w, r, "/produk", http.StatusSeeOther)
}

// Fungsi helper kecil untuk Atoi
func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
