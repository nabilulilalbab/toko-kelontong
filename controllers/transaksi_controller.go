package controllers

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/services"
	"github.com/nabilulilalbab/toko-klontong/utils"
)

type TransaksiController interface {
	ShowKasirPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ProcessCheckout(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ShowHistoryPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ShowHistoryDetailPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GenerateQRIS(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ShowPaymentPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ConfirmPayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
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
	// 1. Decode request dari JSON
	var request services.CreateTransaksiRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Panggil service untuk membuat transaksi (dengan status default 'pending')
	transaksi, err := c.transaksiService.Create(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // Gunakan status 400 untuk error dari client
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// 3. Cek metode pembayaran dan lakukan aksi yang sesuai
	if request.MetodePembayaran == "qris" {
		// Jika QRIS, kirim response JSON yang berisi ID transaksi
		// agar frontend bisa redirect ke halaman pembayaran
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Transaksi dibuat, lanjut ke pembayaran",
			"transaksi_id": transaksi.ID,
			"redirect_url": "/pembayaran/" + strconv.Itoa(int(transaksi.ID)),
		})
	} else { // Asumsikan selain QRIS adalah 'tunai'
		// Jika tunai, langsung update status menjadi 'lunas'
		err := c.transaksiService.UpdateStatus(transaksi.ID, "lunas")
		if err != nil {
			// Handle error jika gagal update status
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Gagal menyelesaikan transaksi tunai"})
			return
		}

		// Kirim response sukses untuk transaksi tunai
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Transaksi tunai berhasil",
			"transaksi_id": transaksi.ID,
		})
	}
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

func (c *transaksiControllerImpl) GenerateQRIS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		StaticQRIS string `json:"static_qris"`
		Total      int    `json:"total"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}

	dynamicQRIS, err := utils.GenerateDynamicQRIS(req.StaticQRIS, uint(req.Total))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"dynamic_qris": dynamicQRIS})
}

func (c *transaksiControllerImpl) ShowPaymentPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	transaksi, err := c.transaksiService.GetByID(uint(id))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Jangan tampilkan QR jika pembayaran bukan QRIS atau sudah lunas
	if transaksi.MetodePembayaran != "qris" || transaksi.Status != "pending" {
		http.Redirect(w, r, "/histori/"+strconv.Itoa(id), http.StatusSeeOther)
		return
	}

	staticPayload := "00020101021126570011ID.DANA.WWW011893600915302259148102090225914810303UMI51440014ID.CO.QRIS.WWW0215ID10200176114730303UMI5204581253033605802ID5922Warung Sayur Bu Sugeng6010Kab. Demak610559567630458C7"
	dynamicPayload, _ := utils.GenerateDynamicQRIS(staticPayload, transaksi.TotalHarga)
	qrCodeBytes, _ := utils.GenerateQRCodeImage(dynamicPayload)
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	data := map[string]interface{}{
		"Transaction": transaksi,
		"QRCode":      qrCodeBase64,
	}
	c.templates.ExecuteTemplate(w, "indexpembayaran.html", data) // Asumsikan nama filenya index.html
}

// 3. TAMBAHKAN ConfirmPayment: untuk update status
func (c *transaksiControllerImpl) ConfirmPayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	c.transaksiService.UpdateStatus(uint(id), "lunas")
	http.Redirect(w, r, "/histori/"+strconv.Itoa(id), http.StatusSeeOther)
}
