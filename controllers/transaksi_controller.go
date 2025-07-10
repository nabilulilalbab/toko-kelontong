package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/services"
)

type TransaksiController interface {
	ShowKasirPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ProcessCheckout(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ShowHistoryPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ShowHistoryDetailPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type transaksiControllerImpl struct {
	transaksiService services.TransaksiService
	produkService    services.ProdukService
	templates        *template.Template
}

func NewTransaksiController(ts services.TransaksiService, ps services.ProdukService, t *template.Template) TransaksiController {
	return &transaksiControllerImpl{
		transaksiService: ts,
		produkService:    ps,
		templates:        t,
	}
}

// Menampilkan halaman utama kasir dengan daftar produk
func (c *transaksiControllerImpl) ShowKasirPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	produks, err := c.produkService.FindAll()
	if err != nil {
		http.Error(w, "Gagal mengambil data produk", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":    "Kasir",
		"Products": produks,
	}

	c.templates.ExecuteTemplate(w, "index.html", data)
}

// Menerima data dari frontend dan memproses checkout
func (c *transaksiControllerImpl) ProcessCheckout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request services.CreateTransaksiRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}

	transaksi, err := c.transaksiService.Create(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Transaksi berhasil",
		"transaksi_id": transaksi.ID,
	})
}

func (c *transaksiControllerImpl) ShowHistoryPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	transaksis, _ := c.transaksiService.GetAll()
	data := map[string]interface{}{
		"Title":        "Riwayat Transaksi",
		"Transactions": transaksis,
	}
	c.templates.ExecuteTemplate(w, "indexhistory.html", data)
}

func (c *transaksiControllerImpl) ShowHistoryDetailPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	transaksi, _ := c.transaksiService.GetByID(uint(id))
	data := map[string]interface{}{
		"Title":       "Detail Transaksi",
		"Transaction": transaksi,
	}
	c.templates.ExecuteTemplate(w, "detailhistory.html", data)
}
