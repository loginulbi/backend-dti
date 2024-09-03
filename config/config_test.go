package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBPresensi(t *testing.T) {
	// Nama database yang ingin diuji
	dbname := "hris"

	// Panggil fungsi DBPresensi
	db := DBPresensi(dbname)

	// Pastikan bahwa koneksi ke database berhasil
	assert.NotNil(t, db, "Expected MongoDB database to be non-nil")

	// Uji koneksi dengan menjalankan ping
	err := db.Client().Ping(nil, nil)
	assert.NoError(t, err, "Expected MongoDB ping to succeed")
}