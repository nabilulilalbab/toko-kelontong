package utils

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/nabilulilalbab/toko-klontong/models"
)

func GenerateProdukExcel(produks []models.Produk) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Laporan Produk"
	f.NewSheet(sheetName)

	// Set header tabel
	headers := []string{"ID", "Nama Produk", "Harga", "Stok"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Isi data produk
	for i, produk := range produks {
		row := i + 2 // Mulai dari baris ke-2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), produk.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), produk.NamaProduk)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), produk.Harga)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), produk.Stok)
	}

	f.DeleteSheet("Sheet1") // Hapus sheet default

	// Tulis file ke buffer di memori
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
