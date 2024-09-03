package atdb

import (
	// "context"
	// "login-service/helper/atdb"
	// "login-service/model"
	"os"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// func TestGetOneDoc(t *testing.T) {
// 	// Connection string ke MongoDB
// 	clientOptions := options.Client().ApplyURI("mongodb://root:MongoDBypbpi123!%40@10.14.200.17:42117")

// 	// Coba sambungkan ke MongoDB
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}

// 	// Cek koneksi dengan mencoba ping MongoDB
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		t.Fatalf("Failed to ping MongoDB: %v", err)
// 	}

// 	// Pilih database dan koleksi
// 	db := client.Database("hris") // Ganti dengan nama database Anda
// 	collection := "karyawan"        // Ganti dengan nama koleksi Anda

// 	// Cari satu dokumen yang sudah ada di database untuk pengujian
// 	filter := bson.M{"nama": "Valen Rionald"} // Pastikan nama ini ada di database Anda

// 	var karyawan model.Karyawan
// 	karyawan, err = GetOneDoc[model.Karyawan](db, collection, filter)

// 	if err != nil {
// 		t.Fatalf("Failed to find document: %v", err)
// 	}

// 	// Verifikasi bahwa dokumen yang ditemukan memiliki field yang diharapkan
// 	assert.Equal(t, "Valen Rionald", karyawan.Nama)
// 	assert.NotEmpty(t, karyawan.PhoneNumber)
// 	assert.NotEmpty(t, karyawan.Email)
// 	assert.NotEmpty(t, karyawan.Jabatan)
// }

// func TestGetOneDoc(t *testing.T) {
//     mconn := DBInfo{
//         DBString: "mongodb://root:MongoDBypbpi123!%40@10.14.200.17:42117",
//         DBName:   "hris",
//     }

//     db, err := MongoConnect(mconn)
//     if err != nil {
//         t.Fatalf("Failed to connect to MongoDB: %v", err)
//     }

//     // Cari satu dokumen yang sudah ada di database untuk pengujian
//     filter := bson.M{"nama": "Valen Rionald"} // Pastikan nama ini ada di database Anda

//     var karyawan model.Karyawan
//     karyawan, err = GetOneDoc[model.Karyawan](db, "karyawan", filter)

//     if err != nil {
//         t.Fatalf("Failed to find document: %v", err)
//     }

//     // Verifikasi bahwa dokumen yang ditemukan memiliki field yang diharapkan
//     assert.Equal(t, "Valen Rionald", karyawan.Nama)
//     assert.NotEmpty(t, karyawan.PhoneNumber)
//     assert.NotEmpty(t, karyawan.Email)
//     assert.NotEmpty(t, karyawan.Jabatan)
// }

func TestMongoConnect(t *testing.T) {
	// Definisikan informasi koneksi MongoDB
	mconn := DBInfo{
		DBString: os.Getenv("MONGOSTRING"), // Ganti sesuai dengan MONGOSTRING Anda
		DBName:   "hris",
	}

	// Panggil fungsi MongoConnect
	db, err := MongoConnect(mconn)

	// Verifikasi bahwa tidak ada error
	assert.NoError(t, err, "Expected no error when connecting to MongoDB")

	// Verifikasi bahwa koneksi tidak nil
	assert.NotNil(t, db, "Expected MongoDB database to be non-nil")

	// Uji koneksi dengan menjalankan ping
	err = db.Client().Ping(nil, nil)
	assert.NoError(t, err, "Expected MongoDB ping to succeed")
}