package errors // Package standar yang kita buat untuk menyeragamkan format error aplikasi

// AppError adalah struktur error kostum buatan kita.
// Keuntungannya menggunakan Custom Error adalah kita bisa menyisipkan Kode dan Pesan Bahasa Manusia
// di samping pesan technical error (yang seringkali berbentuk kode memalukan SQL).
type AppError struct {
	Code    string // Kode error untuk kebutuhan Front-End, misalnya "USER_NOT_FOUND"
	Message string // Pesan bahasa manusia yang aman dikembalikan ke pengguna (Client Safe)
	Err     error  // Menyimpan error inti (misal: sql.ErrNoRows). Hanya dicetak di Log Terminal untuk programmer
}

// Implementasi fungsi Error() membuat struct AppError milik kita
// valid dikenali sebagai tipe "error" bawaan bahasa Go.
func (e *AppError) Error() string {
	return e.Message
}
