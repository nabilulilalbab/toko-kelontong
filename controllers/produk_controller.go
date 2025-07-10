package controllers

import (
	"bytes"
	"encoding/json"
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
	Export(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	SearchAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
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
	hargaStr := r.PostForm.Get("harga")
	stokStr := r.PostForm.Get("stok")

	// --- MULAI BLOK VALIDASI ---
	if namaProduk == "" || hargaStr == "" || stokStr == "" {
		data := map[string]interface{}{
			"Title":      "Tambah Produk Baru",
			"FormAction": "/produk/tambah",
			"Error":      "Semua kolom wajib diisi.",
			// Kirim kembali data yang sudah diinput agar tidak hilang
			"Product": models.Produk{
				NamaProduk: namaProduk,
				Harga:      uint(mustAtoi(hargaStr)),
				Stok:       uint(mustAtoi(stokStr)),
			},
		}
		c.templates.ExecuteTemplate(w, "form.html", data)
		return
	}
	// --- SELESAI BLOK VALIDASI ---

	harga, _ := strconv.Atoi(hargaStr)
	stok, _ := strconv.Atoi(stokStr)

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

// controllers/produk_controller.go

func (c *produkControllerImpl) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal mem-parsing form", http.StatusBadRequest)
		return
	}

	namaProduk := r.PostForm.Get("nama_produk")
	hargaStr := r.PostForm.Get("harga")
	stokStr := r.PostForm.Get("stok")

	// --- MULAI BLOK VALIDASI ---
	if namaProduk == "" || hargaStr == "" || stokStr == "" {
		data := map[string]interface{}{
			"Title":      "Edit Produk",
			"FormAction": "/produk/update/" + strconv.Itoa(id),
			"Error":      "Semua kolom wajib diisi.",
			"Product": models.Produk{
				ID:         uint(id),
				NamaProduk: namaProduk,
				Harga:      uint(mustAtoi(hargaStr)),
				Stok:       uint(mustAtoi(stokStr)),
			},
		}
		c.templates.ExecuteTemplate(w, "form.html", data)
		return
	}
	// --- SELESAI BLOK VALIDASI ---

	produk := models.Produk{
		NamaProduk: namaProduk,
		Harga:      uint(mustAtoi(hargaStr)),
		Stok:       uint(mustAtoi(stokStr)),
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

func (c *produkControllerImpl) Export(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	format := r.URL.Query().Get("format")
	log.Printf("Controller: Memulai proses ekspor ke format: %s.", format)

	produks, err := c.produkService.FindAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data produk", http.StatusInternalServerError)
		return
	}

	var buffer *bytes.Buffer
	var filename string
	var contentType string

	if format == "pdf" {
		buffer, err = utils.GenerateProdukPDF(produks)
		filename = "laporan_produk.pdf"
		contentType = "application/pdf"
	} else if format == "excel" {
		buffer, err = utils.GenerateProdukExcel(produks)
		filename = "laporan_produk.xlsx"
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	} else {
		http.Error(w, "Format tidak didukung", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Gagal membuat file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))
	w.Write(buffer.Bytes())
}

func (c *produkControllerImpl) SearchAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query().Get("q") // Ambil query dari ?q=...
	produks, err := c.produkService.Search(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set header dan kirim response sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produks)
}
