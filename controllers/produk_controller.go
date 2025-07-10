package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/models"
	"github.com/nabilulilalbab/toko-klontong/services"
	"github.com/nabilulilalbab/toko-klontong/utils"
)

// controllers new start

type ProdukController interface {
	Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Form(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ExportExcel(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
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

	produk := models.Produk{
		NamaProduk: namaProduk,
		Harga:      uint(mustAtoi(r.PostForm.Get("harga"))),
		Stok:       uint(mustAtoi(r.PostForm.Get("stok"))),
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
	id := mustAtoi(ps.ByName("id"))
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
	id := mustAtoi(ps.ByName("id"))
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

func (c *produkControllerImpl) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := mustAtoi(ps.ByName("id"))
	log.Printf("Controller: Memulai proses delete produk ID: %d", id)
	err := c.produkService.Delete(uint(id))
	if err != nil {
		log.Printf("Controller: Error saat menghapus produk: %v", err)
		http.Error(w, "Gagal menghapus produk", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/produk", http.StatusSeeOther)
}

func (c *produkControllerImpl) ExportExcel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Controller: Memulai proses ekspor ke Excel.")

	// 1. Ambil semua data produk
	produks, err := c.produkService.FindAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data produk", http.StatusInternalServerError)
		return
	}

	// 2. Panggil generator Excel dari utils
	buffer, err := utils.GenerateProdukExcel(produks)
	if err != nil {
		http.Error(w, "Gagal membuat file Excel", http.StatusInternalServerError)
		return
	}

	// 3. Set HTTP Header untuk memberitahu browser agar men-download file
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_produk.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))

	// 4. Tulis buffer (isi file) ke response
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		http.Error(w, "Gagal mengirim file", http.StatusInternalServerError)
	}
}
