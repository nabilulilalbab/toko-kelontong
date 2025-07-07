package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/toko-klontong/services"
)

var produkTemplate = template.Must(template.ParseFiles("../templates/produk/list.html"))

func ProdukIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	produk := services.GetProdukList()
	err := produkTemplate.Execute(w, produk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
