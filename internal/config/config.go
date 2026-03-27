package config // Package config berfungsi untuk menyimpan struktur konfigurasi environment dari seluruh app

import "time" // Digunakan untuk mengenali tipe durasi seperti detik (s), menit (m)

// Config adalah struktur utama yang membungkus semua konfigurasi aplikasi.
// `mapstructure` adalah tag viper yang dipakai untuk mencocokkan hasil parse dari yaml atau .env ke variable struct ini.
type Config struct {
	App   AppConfig   `mapstructure:"app"`   // Konfigurasi identitas aplikasi
	HTTP  HTTPConfig  `mapstructure:"http"`  // Konfigurasi server HTTP (port & timeout)
	Log   LogConfig   `mapstructure:"log"`   // Konfigurasi format & level log
	Slack SlackConfig `mapstructure:"slack"` // Konfigurasi integrasi webhook ke Slack
	Redis RedisConfig `mapstructure:"redis"` // Konfigurasi server Redis
	DB    DBConfig    `mapstructure:"db"`    // Konfigurasi koneksi database MySQL/PostgreSQL
	JWT   JWTConfig   `mapstructure:"jwt"`   // Konfigurasi rahasia autentikasi token JWT
}

// AppConfig mengatur informasi dasar aplikasi
type AppConfig struct {
	Name string `mapstructure:"name"` // Nama aplikasi kita
	Env  string `mapstructure:"env"`  // Mode environment (contoh: development atau production)
}

// HTTPConfig mengatur prilaku server HTTP saat melayani request
type HTTPConfig struct {
	Port            string        `mapstructure:"port"`             // Port lokal tempat server berjalan
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // Batas waktu membaca koneksi request dari klien
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // Batas waktu untuk menulis response ke klien
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // Batas toleransi durasi saat server dimatikan (graceful shutdown)
}

// LogConfig mengatur seberapa detil log yang ingin kita tampilkan
type LogConfig struct {
	Level string `mapstructure:"level"` // Target log (misalnya: "info", "warn", "error")
}

// SlackConfig mengatur data yang harus dikirimkan log ke integrasi notifikasi Slack
type SlackConfig struct {
	WebhookURL string `mapstructure:"webhook_url"` // URL Hook dari Slack tempat pesan dikirimkan
	Level      string `mapstructure:"level"`       // Batas minimal level log yang akan otomatis dikirim (misal jika "error", hanya error yang dikirim)
	Enabled    bool   `mapstructure:"enabled"`     // Tombol nyala/mati fitur pengiriman log slack
}

// RedisConfig berisi data kredensial akses Redis
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`     // Host dan port server Redis
	Password string `mapstructure:"password"` // Kata sandi jika redis di password (kosong untuk lokal)
	DB       int    `mapstructure:"db"`       // Indeks database Redis (default 0)
}

// DBConfig mengatur grup database.
// Di sistem modular terukur, biasanya kita punya pool koneksi yang beda untuk Read (baca) dan Write (tulis).
type DBConfig struct {
	MySQL    DatabaseGroup `mapstructure:"mysql"`    // Memuat grup read dan write khusus server MySQL
	Postgres DatabaseGroup `mapstructure:"postgres"` // Opsi jika kelak ada PostgreSQL
}

// DatabaseGroup digunakan untuk memisahkan koneksi Master (Write) dan Replica/Slave (Read)
type DatabaseGroup struct {
	Write DatabaseConn `mapstructure:"write"` // Database utama yang digunakan untuk Query INSERT, UPDATE, DELETE
	Read  DatabaseConn `mapstructure:"read"`  // Database tambahan/replika yang dipakai khusus untuk Query SELECT (agar beban database master ringan)
}

// DatabaseConn mendefinisikan URL/DSN (Data Source Name) dan parameter tuning performa per satu titik koneksi
type DatabaseConn struct {
	DSN             string        `mapstructure:"dsn"`                // String URL otentikasi MySQL, e.g. "user:password@tcp(...)"
	MaxOpenConns    int           `mapstructure:"max_open_conns"`     // Limit total maksimal koneksi serentak yang bisa dibuat aplikasi ke DB
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`     // Limit jumlah koneksi yang dibiarkan "hidup stand-by" tidak dipakai
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`  // Batas waktu maksimal satu koneksi boleh didaur ulang, menghindari koneksi busuk/basi
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"` // Lama waktu koneksi nganggur (idle) boleh hidup di atas memori sebelum diputus
}

// JWTConfig memegang kunci rahasia pembuka token dan batas waktu kadaluwarsa token.
type JWTConfig struct {
	AccessSecret    string        `mapstructure:"access_secret"`     // Kata sandi rahasia untuk Token Akses (sebaiknya > 32 karakter acak)
	RefreshSecret   string        `mapstructure:"refresh_secret"`    // Kata sandi rahasia untuk Token Penyegar
	Issuer          string        `mapstructure:"issuer"`            // Pihak yang menerbitkan token (misal: "yourapp-api")
	Audience        string        `mapstructure:"audience"`          // Target pembaca token (misal: "yourapp-client")
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`  // Usia token akses sebelum usang (misal: 15m)
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"` // Usia token penyegar sebelum benar-benar dipaksa login ulang (misal: 168h atau 7 hari)
}
