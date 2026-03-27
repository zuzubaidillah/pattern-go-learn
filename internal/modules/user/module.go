// Package user dipanggil oleh Bootstrapper untuk merakit alat dan modul internal.
// Mengurangi tumpukan kode raksasa berlemak pada main.go.
package user

import (
	"github.com/gin-gonic/gin" // Import router

	"yourapp/internal/infra/db"
	redisx "yourapp/internal/infra/redis"
	"yourapp/internal/modules/user/repository"
	"yourapp/internal/modules/user/service"
	userhttp "yourapp/internal/modules/user/transport/http"
)

// Module adalah abstraksi pengemasan yang mempermudah proses impor dan menyembunyikan service-service
// internal dan repo mentah agar tidak terpapar langsung ke dalam `main.go`.
type Module struct {
	Handler *userhttp.Handler // Handler HTTP saja yang digendong dan diekspos keluar ruangan karena bootstrap hanya butuh Handle Routes
}

// NewModule menyuntikkan ketergantungan (Dependency Injection) hierarkis milik pendaftar (User).
// Semua alur dan alat komunikasi Module ditenun dan disambung pipanya melalui fungsi instansiasi modul ini.
func NewModule(dbm *db.Manager, cache *redisx.Cache) *Module {
	// Layer Repositori dibuat paling mentah: Sambungkan ke Database Read/Write.
	repo := repository.NewSQLRepository(dbm.MySQLWrite, dbm.MySQLRead)

	// Layer Servis menggunakan Repositori tadi sekaligus ditambah fitur Caching.
	svc := service.New(repo, cache)

	// Layer Handler memakai Servis yang sudah dibumbui dengan Cache.
	h := userhttp.NewHandler(svc)

	// Kembalikan kotak bundelan modul untuk dikonsumsi API Framework eksternal
	return &Module{
		Handler: h,
	}
}

// RegisterRoutes membungkus ulang (adapter) untuk fungsi rute HTTP kita
// Menjadikan alur impor di `app.go` lebih bersih karena mereka tidak perlu merogoh
// folder jauh ke dalam hingga `user/transport/http.RegisterRoutes` setiap kali rilis / perubahan struktur folder.
func RegisterRoutes(rg *gin.RouterGroup, h *userhttp.Handler) {
	userhttp.RegisterRoutes(rg, h)
}
