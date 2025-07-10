package utils

import "time"

// FormatTanggal akan mengubah format waktu sesuai standar yang kita mau
func FormatTanggal(t time.Time) string {
	// Format "dd-mm-yyyy hh:mm:ss"
	return t.Format("02 Jan 2006 15:04:05")
}
