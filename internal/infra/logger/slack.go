package logger // Tempat berkumpulnya semua tool catatan sistem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore" // Fundamental core logger dari driver zap
	"yourapp/internal/config"
)

// parseLevel adalah helper kecil guna mengkonversi string ("info", "error") dari yaml
// menjadi angka konstanta Enum zapcore.Level.
func parseLevel(l string) zapcore.Level {
	var level zapcore.Level
	// Menjalankan unmarshal string ke tipe bawaan Zap
	if err := level.UnmarshalText([]byte(l)); err != nil {
		return zapcore.InfoLevel // Secara standar jika gagal, jatuhkan ke tingkatan Info
	}
	return level
}

// SlackHook mengembalikan sebuah fungsi berantai (hook) yang akan dicegat secara otomatis oleh Zap
// saat mencetak Log, lalu data log-nya dikirim secara HTTP POST ke webhook Slack.
func SlackHook(cfg config.SlackConfig) func(zapcore.Entry) error {
	// Membuat agen HTTP penembak dengan timeout 5 detik untuk mencegah hang kalau Slack melambat
	client := &http.Client{Timeout: 5 * time.Second}

	// Tarik standar minimum tingkatan dari konfigurasi
	targetLevel := parseLevel(cfg.Level)

	// Kembalikan closure fungsi yang ditangkap oleh zapcore
	return func(entry zapcore.Entry) error {
		// Validasi filterting: Hanya tembak log ke Slack apabila tingkat log saat ini
		// (entry.Level) lebih gawat atau setara dari level minimal Slack yang di yaml
		if entry.Level < targetLevel {
			return nil
		}

		// Merangkai isi chat slack (contoh: "[ERROR] Gagal membaca Redis")
		payload := map[string]string{
			"text": fmt.Sprintf("[%s] %s", entry.Level.CapitalString(), entry.Message),
		}

		// Ubah jadi JSON string dan tembak permintaannya ke channel slack via HTTP REST
		body, _ := json.Marshal(payload)
		resp, err := client.Post(cfg.WebhookURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		// Bersihkan memory dan matikan alur request HTTP setelah membalas Slack
		defer resp.Body.Close()
		return nil
	}
}
