package logger

import (
	"go.uber.org/zap" // Library logging dari Uber. Jauh lebih kencang dibanding log bawaan Go.
	"yourapp/internal/config"
)

// New menciptakan inti sistem perekam informasi/catatan aplikasi.
func New(cfg config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	// Menghidupkan mode produksi jika "env: production"
	// Mode Produksi keluaran lognya tidak berwarna, super padat berformat JSON murni untuk mesin
	if cfg.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		// Mode Development lognya berwarna warni, rapi baris per baris, gampang dibaca programer
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err // Jika gagal rakit logger
	}

	// Menggabungkan Hook slack buatan kita ke dalam rantai kehidupan objek Zap logger
	if cfg.Slack.Enabled && cfg.Slack.WebhookURL != "" {
		// logger.WithOptions menciptakan duplikat dari logger, lalu menambahkan opsi "selalu panggil fungsi SlackHook" jika print log terjadi.
		logger = logger.WithOptions(zap.Hooks(SlackHook(cfg.Slack)))
	}

	return logger, nil
}
