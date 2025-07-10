package utils

import (
	"fmt"
	"strings"

	"github.com/sigurn/crc16"
	"github.com/skip2/go-qrcode"
)

// GenerateDynamicQRIS mengonversi payload QRIS statis menjadi dinamis.
func GenerateDynamicQRIS(staticPayload string, amount uint) (string, error) {
	payload := staticPayload[:len(staticPayload)-4]
	payload = strings.Replace(payload, "010211", "010212", 1)
	amountStr := fmt.Sprintf("%d", amount)
	transactionAmount := fmt.Sprintf("54%02d%s", len(amountStr), amountStr)
	parts := strings.Split(payload, "5802ID")
	payload = parts[0] + transactionAmount + "5802ID" + parts[1]

	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(payload), table)
	crcStr := fmt.Sprintf("%04X", crc)

	return payload + crcStr, nil
}

// GenerateQRCodeImage membuat gambar QR code dari payload dan mengembalikannya sebagai byte PNG.
func GenerateQRCodeImage(payload string) ([]byte, error) {
	png, err := qrcode.Encode(payload, qrcode.Medium, 256) // 256x256 pixels
	if err != nil {
		return nil, err
	}
	return png, nil
}

