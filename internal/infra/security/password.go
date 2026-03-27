package security // Package security menyimpan utilitas fundamental otentikasi. Termasuk Hash sandi & Enkripsi data.

import "golang.org/x/crypto/bcrypt" // Library enkripsi standar dari golang

// BcryptHasher melambangkan algoritma pengacak sandi untuk implementasi Auth
type BcryptHasher struct{}

// NewBcryptHasher pendiri paten alat keamanan Bcrypt baru
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Verify mengadu silang / komparasi sandi asli ketikan log in (plain) versus sandi hash dalam database terenkripsi panjang.
// Ia merupakan implementasi wujud konkret dari interface janji suci `authservice.PasswordVerifier` nantinya.
func (b *BcryptHasher) Verify(hashedPassword, plainPassword string) bool {
	// Membandingkan isi. Karena bcrypt mensyaratkan list []byte, kita konversi string ke byte mentah terlebih dulu.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	// Kalau error bernilai `nil`, tandanya perbandingan password hash lolos 100%. User memasukkan sandi aslinya.
	return err == nil
}

// Hash mengubah text rahasia murni (misal '123456') menjadi '$$2a$10..xxxzy'
// agar tak ada teknisi/hacker yang dapat membacanya langsung via kolom database.
//
// Metode Hash ini sengaja disediakan meski tak dipakai login (karena untuk kelengkapan Registrasi module yang belum diimplement).
func (b *BcryptHasher) Hash(plainPassword string) (string, error) {
	// Cost 10 (DefaultCost) adalah limit aman komputasi Bcrypt tanpa memberhalangi resourcess memori CPU server Go Anda.
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
