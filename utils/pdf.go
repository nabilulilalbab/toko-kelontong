package utils

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"github.com/nabilulilalbab/toko-klontong/models"
)

func GenerateProdukPDF(produks []models.Produk) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Judul Dokumen
	pdf.Cell(40, 10, "Laporan Data Produk")
	pdf.Ln(20) // Spasi ke bawah

	// Header Tabel
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(20, 10, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, 10, "Nama Produk", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Harga", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Stok", "1", 0, "C", false, 0, "")
	pdf.Ln(10)

	// Isi Tabel
	pdf.SetFont("Arial", "", 12)
	for _, p := range produks {
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", p.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, 10, p.NamaProduk, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", p.Harga), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", p.Stok), "1", 0, "R", false, 0, "")
		pdf.Ln(10)
	}

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	return &buffer, nil
}
