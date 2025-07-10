package utils

import (
	"html/template"
	"log"
)

func ParseTemplates(basePath string) *template.Template {
	log.Println("Parsing semua templates...")
	funcMap := template.FuncMap{
		"formatTanggal": FormatTanggal,
	}
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(basePath + "templates/**/*.html")
	if err != nil {
		// Jika ada error (misal: pola salah), program akan berhenti dengan pesan jelas
		panic("Gagal mem-parsing templates: " + err.Error())
	}

	log.Println("Parsing templates selesai.")
	return tmpl
}

