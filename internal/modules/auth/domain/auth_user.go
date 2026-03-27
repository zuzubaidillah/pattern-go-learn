package domain

// AuthUser beda halnya dengan entitas modul User utama!
// Model ini secara sepihak dan independen menampung kolom terpenting login saja (seperti Sandi & Status).
// Di arsitektur Modular, Modul Auth TAK BOLEH langsung saling panggil entitas modul User lain tanpa pembatas,
// agar tetap rapi jika project dipisahkan (Microservice-ready).
type AuthUser struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string
	Status       string
}
