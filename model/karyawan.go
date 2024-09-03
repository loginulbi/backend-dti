package model

type JamKerja struct {
	Durasi    int      `bson:"durasi" json:"durasi"`
	JamMasuk  string   `bson:"jam_masuk" json:"jam_masuk"`
	JamKeluar string   `bson:"jam_keluar" json:"jam_keluar"`
	Gmt       int      `bson:"gmt" json:"gmt"`
	Hari      []string `bson:"hari" json:"hari"`
}

type Karyawan struct {
	Nama        string    `bson:"nama" json:"nama"`
	PhoneNumber string    `bson:"phone_number" json:"phone_number"`
	Email       string    `bson:"email" json:"email"`
	Jabatan     string    `bson:"jabatan" json:"jabatan"`
	JamKerja    []JamKerja `bson:"jam_kerja" json:"jam_kerja"`
	HariKerja   []string  `bson:"hari_kerja" json:"hari_kerja"`
}
