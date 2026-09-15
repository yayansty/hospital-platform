package models

type RoomAvailability struct {
	KodeRuang string `json:"koderuang"`
	NamaRuang string `json:"nama_ruang"`
	Kapasitas int64  `json:"kapasitas"`
	Isi       int64  `json:"isi"`
	Disi      int64  `json:"disi"`
	Urut      int64  `json:"urutan"`
	KodeKelas string `json:"kodekelas"`
	NamaKelas string `json:"nama_kelas"`
	Tersedia  int64  `json:"tersedia"`
}

