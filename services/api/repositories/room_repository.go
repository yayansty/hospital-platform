package repositories

import (
	"database/sql"

	"medic-api/helpers"
	"medic-api/models"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) FindAvailability() ([]models.RoomAvailability, error) {
	query := `
		SELECT
			A.KODERUANG,
			A.NAMA_RUANG,
			A.KAPASITAS,
			A.ISI,
			A.DISI,
			A.URUT,
			A.KODEKELAS,
			A.NAMA_KELAS,
			CASE
				WHEN A.TERSEDIA < 0 THEN 0
				ELSE A.TERSEDIA
			END AS TERSEDIA
		FROM (
			SELECT
				A.KODERUANG,
				A.NAMA_RUANG,
				A.KAPASITAS,
				A.ISI,
				A.DISI,
				A.URUT,
				A.KODEKELAS,
				A.NAMA_KELAS,
				A.KAPASITAS - A.DISI - A.ISI AS TERSEDIA
			FROM (
				SELECT
					A.KODERUANG,
					A.NAMA_RUANG,
					A.KAPASITAS,
					A.ISI,
					A.URUT,
					A.KODEKELAS,
					A.NAMA_KELAS,
					SUM(A.DISI) AS DISI
				FROM (
					SELECT
						A.KODERUANG,
						C.NAMA_RUANG,
						C.KAPASITAS,
						C.ISI,
						C.URUT,
						A.KODEKELAS,
						B.NAMA_KELAS,
						(A.JUMLAH_BED + A.BEDD - A.BOOKING) AS JUMLAH_BED,
						(
							SELECT COUNT(*)
							FROM USER_TMC.REGISTRASI
							WHERE ID_UNIT = '100'
								AND ID_KAMAR = A.ID_KAMAR
								AND KELUAR = 0
								AND TGL_END IS NULL
						) AS DISI
					FROM USER_TMC.KAMAR A
					LEFT JOIN USER_TMC.BPJS_KELAS B
						ON B.ID_KELAS = A.KODEKELAS
					LEFT JOIN USER_TMC.BPJS_RUANG C
						ON C.KODERUANG = A.KODERUANG
					WHERE A.BPJS = 1
						AND A.KODEKELAS IS NOT NULL
						AND A.KODERUANG IS NOT NULL
						AND C.KAPASITAS > 0
				) A
				GROUP BY
					A.KODERUANG,
					A.NAMA_RUANG,
					A.KAPASITAS,
					A.ISI,
					A.URUT,
					A.KODEKELAS,
					A.NAMA_KELAS
			) A
		) A
		ORDER BY A.URUT
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, helpers.ErrDatabase
	}
	defer rows.Close()

	var rooms []models.RoomAvailability

	for rows.Next() {
		var room models.RoomAvailability

		err := rows.Scan(
			&room.KodeRuang,
			&room.NamaRuang,
			&room.Kapasitas,
			&room.Isi,
			&room.Disi,
			&room.Urut,
			&room.KodeKelas,
			&room.NamaKelas,
			&room.Tersedia,
		)

		if err != nil {
			return nil, helpers.ErrDatabase
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, helpers.ErrDatabase
	}

	return rooms, nil
}
