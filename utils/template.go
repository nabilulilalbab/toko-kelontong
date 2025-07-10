package utils

import (
	"html/template"
	"log"
)

func ParseTemplates() *template.Template {
	log.Println("Parsing semua templates...")
	funcMap := template.FuncMap{
		"formatTanggal": FormatTanggal,
	}
	// Gunakan "../" untuk naik satu level dari direktori 'cmd/'
	// Gunakan "**/" untuk mencari di semua sub-folder secara rekursif
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/**/*.html")
	if err != nil {
		// Jika ada error (misal: pola salah), program akan berhenti dengan pesan jelas
		panic("Gagal mem-parsing templates: " + err.Error())
	}

	log.Println("Parsing templates selesai.")
	return tmpl
}

