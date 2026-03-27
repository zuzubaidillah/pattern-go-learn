package bootstrap // Package bootstrap adalah tempat semua logika inisialisasi awal program berkumpul

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin" // Gin adalah framework web andalan kita untuk routing HTTP
	"github.com/spf13/viper"   // Viper membantu meng-import konfigurasi dari .env atau yaml ke struct go
	"go.uber.org/zap"          // Zap untuk sistem logging yang efisien dan detail

	"yourapp/internal/config"             // Mengambil tipe-tipe struct konfigurasi yang kita buat
	"yourapp/internal/infra/db"           // Infrastruktur database
	"yourapp/internal/infra/logger"       // Infrastruktur logger
	redisx "yourapp/internal/infra/redis" // Infrastruktur redis (dinamai redisx agar tidak tabrakan dengan package aslinya)
	"yourapp/internal/modules/auth"       // Import module Auth untuk registrasi token JWT
	"yourapp/internal/modules/user"       // Import module User yang menyimpan logika tentang pendaftaran pengguna
)

// App adalah struktur pembungkus kerangka layanan utama kita
type App struct {
	Server *http.Server // Node Server HTTP internal Go
}

// NewApp bertugas merakit seluruh aplikasi: mulai dari membaca config, database hingga memasang route HTTP gin
func NewApp() (*App, error) {
	// 1. Membaca setting dari environment variables / yaml ke memori
	cfg := MustLoadConfig()

	// 2. Menginisialisasi logger dengan setting yang telah divariabelisasi (termasuk fitur Slack)
	logg, err := logger.New(cfg)
	if err != nil {
		return nil, err
	}

	// 3. Merakit koneksi database (MySQL Read / Write Pools)
	dbm, err := MustBuildDBManager(cfg)
	if err != nil {
		// Log yang sangat terstruktur: "Mencetak pesan Fatal" + lampiran payload error.
		logg.Fatal("Failed to coordinate DB connections", zap.Error(err))
		return nil, err
	}
	// Khusus untuk DB, kita buat peringatan informasi bila koneksinya berhasil tersambung
	logg.Info("DB Listen: Database connections established successfully")

	// 4. Menyambungkan koneksi Redis untuk caching
	rdb, err := redisx.NewClient(cfg.Redis)
	if err != nil {
		return nil, err
	}
	// Bungkus klien redis mentah ke dalam helper `Cache` kita
	cache := redisx.NewCache(rdb)

	// Menginstansiasi loket penyimpanan Refresh Token Sesi Jangka Panjang.
	refreshStore := redisx.NewRefreshStore(rdb)

	// 5. Membuat kerangka sistem routing web dengan framework Gin
	r := gin.New()
	r.Use(gin.Recovery()) // Middleware agar kalau server crash / panic, akan merecovery otomatis (tidak down server)
	r.Use(gin.Logger())   // Middleware standar penyematan log info untuk setiap API Call dari Gin

	// Mengelompokkan semua rute di bawah "namadomain.com/api/v1/"
	api := r.Group("/api/v1")

	// 6. Menyuntikkan infrastruktur (Database, Cache) ke "Modul User".
	// DI (Dependency Injection) membuat layanan user tidak berdiri sendiri tapi diberikan resources lewat parameternya
	userModule := user.NewModule(dbm, cache)

	// Menyuntikkan alat yang sama dan alat ekstra (refreshStore) ke "Modul Auth".
	authModule := auth.NewModule(cfg, dbm, cache, refreshStore)

	// Daftarkan jalur (URL PATH) User dan Auth ke pelayan rute API `v1`.
	user.RegisterRoutes(api, userModule.Handler)
	auth.RegisterRoutes(api, authModule)

	// 7. Jadikan instansi router `r` tadi handler dari server HTTP Go mentah. Ini memastikan port diset sesuai config.
	server := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: r,
	}

	// Kembalikan App yang sudah siap diluncurkan
	return &App{Server: server}, nil
}

// Run adalah trigger untuk memulai blokade "mendengar dari Port"
func (a *App) Run() error {
	log.Printf("Starting server on %s", a.Server.Addr)
	return a.Server.ListenAndServe() // Perintah blok aktif agar server tidak mati sendirinya
}

// MustLoadConfig memastikan membaca config. Jika gagal membaca konfigurasi, matikan program (Must = harus bisa)
func MustLoadConfig() config.Config {
	viper.AddConfigPath("./configs")      // Cari di folder `configs`
	viper.SetConfigName("config.example") // Cari file spesifik bernama config.example
	viper.SetConfigType("yaml")           // Pakai ekstensi .yaml

	// Kalau ada nested variable yaml (contoh app.name) gantikan titik jadi underscore jika dibaca di ENV OS lokal, misal APP_NAME
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // Tumpang tindihkan file yaml dengan file dari .env / env sistem OS bila terdeteksi nama yang sama

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg config.Config
	// Konversikan data hasil baca ke memori Object (struct) Go
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return cfg
}

// MustBuildDBManager merakit dua kolam MySQL Read dan Write yang kita set up di ENV
func MustBuildDBManager(cfg config.Config) (*db.Manager, error) {
	// Menghidupkan kolam ke database write
	writePool, err := db.NewMySQL(cfg.DB.MySQL.Write)
	if err != nil {
		return nil, err
	}

	// Menghidupkan kolam ke database read
	readPool, err := db.NewMySQL(cfg.DB.MySQL.Read)
	if err != nil {
		return nil, err
	}

	// Menyimpannya ke satu Manajer utama agar gampang di Passing
	return &db.Manager{
		MySQLWrite: writePool,
		MySQLRead:  readPool,
	}, nil
}
