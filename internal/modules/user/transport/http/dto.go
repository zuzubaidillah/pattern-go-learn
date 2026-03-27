package http

// CreateUserRequest adalah DTO (Data Transfer Object).
// Berbeda dari Entity Domain yang abstrak, Struct ini khusus melayani Validasi JSON dan Format Parsing.
type CreateUserRequest struct {
	// Tag JSON memastikan key API Body-nya `"name": ""`.
	// Tag binding akan menolak permintaan otomatis dengan kode 400 Bad Request jika nama ini < 2 abjad atau kosong.
	Name string `json:"name" binding:"required,min=2,max=100"`

	// Validasi regex tipe surel agar klien tak asal main isi ID tak berformat!
	Email string `json:"email" binding:"required,email"`

	// Enum kera-kera di back-end: Hanya menoleransi 'active' atau 'inactive'
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

// UpdateUserRequest meregulasi perubahan payload (Payload Schema).
type UpdateUserRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=100"`
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

// UserResponse adalah model tampilan JSON final ke Client Postman.
// Kita mencegah kata sandi atau salt field ikut terkirim dengan menyaringnya di struct ini.
type UserResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}
