package repository

import (
	"database/sql"
	"errors"
	"strings"
)

type Mahasiswa struct {
	ID                 int      `json:"id"`
	Nama               string   `json:"nama"`
	Usia               int      `json:"usia"`
	Gender             string   `json:"gender"`
	Tanggal_Registrasi string   `json:"tanggal_registrasi"`
	ID_Jurusan         string   `json:"id_jurusan"`
	Hobi               []string `json:"hobi"`
}

type MahasiswaRepository interface {
	GetMahasiswaById(id int) (*Mahasiswa, error)
	GetAllMahasiswa() ([]Mahasiswa, error)
	InsertMahasiswa(mahasiswa *Mahasiswa) (*Mahasiswa, error)
	UpdateMahasiswa(mahasiswa *Mahasiswa, id int) (*Mahasiswa, error)
	DeleteMahasiswa(id int) error
}

type MahasiswaRepo struct {
	db *sql.DB
}

func NewMahasiswaRepository(db *sql.DB) MahasiswaRepository {
	return &MahasiswaRepo{
		db: db,
	}
}

func (r *MahasiswaRepo) GetMahasiswaById(id int) (*Mahasiswa, error) {
	mahasiswa := Mahasiswa{}

	query := "SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, GROUP_CONCAT(h.nama_hobi SEPARATOR ';') FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id WHERE m.id = ? GROUP BY m.id"
	row := r.db.QueryRow(query, id)

	var hobi string
	err := row.Scan(
		&mahasiswa.ID,
		&mahasiswa.Nama,
		&mahasiswa.Usia,
		&mahasiswa.Gender,
		&mahasiswa.Tanggal_Registrasi,
		&mahasiswa.ID_Jurusan,
		&hobi,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No rows found, so return nil instead of an error
		}
		return nil, err
	}

	// Convert integer Gender value to string
	if mahasiswa.Gender == "0" {
		mahasiswa.Gender = "male"
	} else if mahasiswa.Gender == "1" {
		mahasiswa.Gender = "female"
	}

	// Split comma-separated string of hobbies into a slice of strings
	hobiArr := strings.Split(hobi, ";")
	mahasiswa.Hobi = hobiArr

	return &mahasiswa, nil
}

func (r *MahasiswaRepo) GetAllMahasiswa() ([]Mahasiswa, error) {
	allMahasiswa := []Mahasiswa{}
	query := "SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, GROUP_CONCAT(h.nama_hobi SEPARATOR ';') FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id GROUP BY m.id"

	result, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var hobi string
		var mahasiswa Mahasiswa
		err := result.Scan(&mahasiswa.ID,
			&mahasiswa.Nama,
			&mahasiswa.Usia,
			&mahasiswa.Gender,
			&mahasiswa.Tanggal_Registrasi,
			&mahasiswa.ID_Jurusan,
			&hobi)
		if err != nil {
			return nil, err
		}

		// Convert integer Gender value to string
		if mahasiswa.Gender == "0" {
			mahasiswa.Gender = "male"
		} else if mahasiswa.Gender == "1" {
			mahasiswa.Gender = "female"
		}

		// Split comma-separated string of hobbies into a slice of strings
		hobiArr := strings.Split(hobi, ";")
		mahasiswa.Hobi = hobiArr

		allMahasiswa = append(allMahasiswa, mahasiswa)
	}
	return allMahasiswa, nil
}

func (r *MahasiswaRepo) InsertMahasiswa(mahasiswa *Mahasiswa) (*Mahasiswa, error) {
	insertQueryMahasiswa := "INSERT INTO mahasiswa(nama, usia, gender, tanggal_registrasi, id_jurusan) VALUES (?, ?, ?, ?, ?)"
	insertQueryHobbies := "INSERT INTO mahasiswa_hobi(id_mahasiswa, id_hobi) VALUES (?, ?)"

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(insertQueryMahasiswa)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	result, err := stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, mahasiswa.ID_Jurusan)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, h := range mahasiswa.Hobi {
		stmt, err := tx.Prepare(insertQueryHobbies)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = stmt.Exec(lastID, h)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return mahasiswa, nil
}

func (r *MahasiswaRepo) UpdateMahasiswa(mahasiswa *Mahasiswa, id int) (*Mahasiswa, error) {
	queryUpdate := "UPDATE mahasiswa SET nama = ?, usia = ?, gender = ?, tanggal_registrasi = ?, id_jurusan = ? WHERE id = ?"
	queryDelete := "DELETE FROM mahasiswa_hobi WHERE id_mahasiswa = ?"
	queryInsertHobbies := "INSERT INTO mahasiswa_hobi(id_mahasiswa, id_hobi) VALUES (?, ?)"

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(queryUpdate)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, mahasiswa.ID_Jurusan, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(queryDelete, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, h := range mahasiswa.Hobi {
		stmt, err = tx.Prepare(queryInsertHobbies)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = stmt.Exec(id, h)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return mahasiswa, nil
}

func (r *MahasiswaRepo) DeleteMahasiswa(id int) error {
	isExistQuery := "SELECT id FROM mahasiswa WHERE id = ?"
	deleteHobbiesQuery := "DELETE FROM mahasiswa_hobi WHERE id_mahasiswa =?"
	deleteMahasiswaQuery := "DELETE FROM mahasiswa WHERE id =?"
	// Transaction declaration
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if mahasiswa exists in database
	var mahasiswa Mahasiswa
	err = tx.QueryRow(isExistQuery, id).Scan(&mahasiswa.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	// Delete the corresponding rows in mahasiswa_hobi table
	hobiStmt, err := tx.Prepare(deleteHobbiesQuery)
	if err != nil {
		return err
	}
	_, err = hobiStmt.Exec(id)
	if err != nil {
		return err
	}

	// Delete the mahasiswa
	mahasiswaStmt, err := tx.Prepare(deleteMahasiswaQuery)
	if err != nil {
		return err
	}
	_, err = mahasiswaStmt.Exec(id)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
