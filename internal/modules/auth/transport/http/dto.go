package http // Cangkang pembungkus gerbang awal (Form Handling dari JSON Postman ke Golang Data)

// LoginRequest mewadahi syarat form saat orang mengetik endpoint POST /login.
type LoginRequest struct {
	// Wajib isi, harus lolos regex standar nama domain surel.
	Email string `json:"email" binding:"required,email"`
	// Minimal ketik sandi 6 abjad. Demi meluaskan area keamanan, password dipatok ketat tanpa dibatasi limit karakter maksimum berlebihan.
	Password string `json:"password" binding:"required,min=6"`
}

// RefreshRequest mencomot body pas mendarat di endpoint perpanjangan POST /refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenResponse adalah standardisasi seragam format penyusunan output kembalian Paket Autorisasi agar mematuhi konvensi industri OAuth2.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`  // Karcis Masuk sesi singgah API harian
	RefreshToken string `json:"refresh_token"` // Kunci master kendali siklus perpanjangan Sesi
	TokenType    string `json:"token_type"`    // Lazim diisi patokan konvensi Web: "Bearer"
}
