package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// calculateCRC16 (Checksum) adalah terjemahan dari fungsi ConvertCRC16 di PHP
func calculateCRC16(data string) string {
	crc := 0xFFFF
	for _, b := range []byte(data) {
		crc ^= int(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) > 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}
	// Format ke Hex 4 digit dengan padding nol jika perlu
	return fmt.Sprintf("%04X", crc&0xFFFF)
}

// GenerateDynamicQRIS adalah fungsi utama untuk membuat QRIS dinamis
func GenerateDynamicQRIS(staticQRIS string, amount int) (string, error) {
	if len(staticQRIS) < 4 {
		return "", errors.New("string QRIS statis tidak valid")
	}

	// Hapus 4 karakter checksum lama
	qrisWithoutChecksum := staticQRIS[:len(staticQRIS)-4]

	// Ganti indikator statis (010211) ke dinamis (010212)
	// Hanya ganti kemunculan pertama untuk keamanan
	step1 := strings.Replace(qrisWithoutChecksum, "010211", "010212", 1)

	// Pisahkan string berdasarkan "5802ID" untuk menyisipkan nominal
	parts := strings.Split(step1, "5802ID")
	if len(parts) != 2 {
		return "", errors.New("format QRIS tidak mengandung '5802ID' yang valid")
	}

	// Buat tag nominal (Tag 54)
	amountStr := strconv.Itoa(amount)
	amountTag := fmt.Sprintf("54%02d%s", len(amountStr), amountStr)

	// Gabungkan kembali string dengan tag nominal
	finalQRIS := parts[0] + amountTag + "5802ID" + parts[1]

	// Hitung checksum baru dan tambahkan ke akhir
	newChecksum := calculateCRC16(finalQRIS)
	finalQRIS += newChecksum

	return finalQRIS, nil
}
