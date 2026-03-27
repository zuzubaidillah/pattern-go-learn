package auth // Menyambung kepingan modul ke dunia luar (Faktor Wiring & DI Container)

import (
	"github.com/gin-gonic/gin"

	"yourapp/internal/config"
	"yourapp/internal/infra/db"
	redisx "yourapp/internal/infra/redis"
	"yourapp/internal/infra/security"
	"yourapp/internal/modules/auth/repository"
	authservice "yourapp/internal/modules/auth/service"
	authhttp "yourapp/internal/modules/auth/transport/http"
)

// Module sebagai penyimpan kerangka utuh
type Module struct {
	Handler *authhttp.Handler
	Service *authservice.Service
}

// NewModule menyusun dan merangkai instansi dari paling mentah (Database) ke paling luhur (HTTP Handler)
func NewModule(cfg config.Config, dbm *db.Manager, cache *redisx.Cache, redisStore *redisx.RefreshStore) *Module {

	// Alat Enkripsi Keamanan didirikan
	passVerifier := security.NewBcryptHasher()

	// Gudang data spesifik pembacaan Login ditautkan ke DB Replica
	repo := repository.NewSQLRepository(dbm.MySQLRead)

	// Pembuahan Inti Logic Domain Service (Dilengkapi Secret JWT, Repo, Redis Store, dan Alat Bcrypt)
	svc := authservice.New(cfg.JWT, repo, redisStore, passVerifier)

	// Pembungkus gerbang lalu lintas Web
	h := authhttp.NewHandler(svc)

	return &Module{
		Handler: h,
		Service: svc,
	}
}

// RegisterRoutes menjahit handler yang ditenun `NewModule` ke dalam alur pipa API Utama Gin.
func RegisterRoutes(rg *gin.RouterGroup, m *Module) {
	// Kita passing Service ke middleware agar satpam `AuthJWT` bisa membongkar kunci token dari User Request
	authhttp.RegisterRoutes(rg, m.Handler, m.Service)
}
