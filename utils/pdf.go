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

func GenerateTransaksiPDF(transaksis []models.Transaksi) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Laporan Riwayat Transaksi")
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(20, 10, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Total Harga", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Metode Pembayaran", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Status", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for _, t := range transaksis {
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", t.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", t.TotalHarga), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, t.MetodePembayaran, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, t.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, FormatTanggal(t.CreatedAt), "1", 0, "C", false, 0, "")
		pdf.Ln(10)
	}

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	return &buffer, nil
}
