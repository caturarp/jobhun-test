package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, GROUP_CONCAT(h.nama_hobi) FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id WHERE m.id =? GROUP BY m.id", params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving mahasiswa data"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		panic(err.Error())
	}
	defer result.Close()

	var mahasiswa Mahasiswa
	for result.Next() {
		err := result.Scan(&mahasiswa.ID,
			&mahasiswa.Nama,
			&mahasiswa.Usia,
			&mahasiswa.Gender,
			&mahasiswa.Tanggal_Registrasi,
			&mahasiswa.ID_Jurusan,
			&mahasiswa.Hobi)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error scanning mahasiswa data"))
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			panic(err.Error())
		}

		// Convert integer Gender value to string
		if mahasiswa.Gender == "0" {
			mahasiswa.Gender = "male"
		} else if mahasiswa.Gender == "1" {
			mahasiswa.Gender = "female"
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mahasiswa)
	}

}

func getAllMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var allMahasiswa []Mahasiswa

	result, err := db.Query("SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, GROUP_CONCAT(h.nama_hobi) FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id GROUP BY m.id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving mahasiswa data"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var mahasiswa Mahasiswa
		err := result.Scan(&mahasiswa.ID,
			&mahasiswa.Nama,
			&mahasiswa.Usia,
			&mahasiswa.Gender,
			&mahasiswa.Tanggal_Registrasi,
			&mahasiswa.ID_Jurusan,
			&mahasiswa.Hobi)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error scanning mahasiswa data"))
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			panic(err.Error())
		}
		// Convert integer Gender value to string
		if mahasiswa.Gender == "0" {
			mahasiswa.Gender = "male"
		} else if mahasiswa.Gender == "1" {
			mahasiswa.Gender = "female"
		}
		allMahasiswa = append(allMahasiswa, mahasiswa)
	}
	if len(allMahasiswa) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No mahasiswa data found"))
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allMahasiswa)
}

func insertMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
			return
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO mahasiswa(nama, usia, gender, tanggal_registrasi, id_jurusan) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "failed to prepare insert statement", http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var mahasiswa Mahasiswa
	err = json.Unmarshal(body, &mahasiswa)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	result, err := stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, mahasiswa.ID_Jurusan)
	if err != nil {
		http.Error(w, "failed to execute insert statement ", http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "failed to get last inserted ID", http.StatusInternalServerError)
		return
	}
	for _, h := range mahasiswa.Hobi {
		stmt, err = tx.Prepare("INSERT INTO mahasiswa_hobi(id_mahasiswa, id_hobi) VALUES (?, ?)")
		if err != nil {
			http.Error(w, "failed to prepare insert statement for mahasiswa_hobi table", http.StatusInternalServerError)
			panic(err.Error())
		}

		_, err = stmt.Exec(lastID, h)
		if err != nil {
			http.Error(w, "failed to execute query", http.StatusInternalServerError)
			panic(err.Error())

		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "new mahasiswa was created"})
	fmt.Fprintf(w, "New mahasiswa was created")
}

func updateMahasiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	mahasiswaStmt, err := db.Prepare("UPDATE mahasiswa SET nama = ?, usia = ?, gender =?, tanggal_registrasi =?, id_jurusan = ? WHERE mahasiswa_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer mahasiswaStmt.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mahasiswa Mahasiswa
	err = json.Unmarshal(body, &mahasiswa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = mahasiswaStmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, mahasiswa.ID_Jurusan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// delete existing Hobi
	hobiStmt, err := tx.Prepare("DELETE FROM mahasiswa_hobi WHERE id_mahasiswa=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer hobiStmt.Close()

	_, err = hobiStmt.Exec(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// insert new Hobi
	hobiInsertStmt, err := tx.Prepare("INSERT INTO mahasiswa_hobi(id_mahasiswa, id_hobi) VALUES (?,?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer hobiInsertStmt.Close()

	for _, hobiID := range mahasiswa.Hobi {
		_, err = hobiInsertStmt.Exec(params["id"], hobiID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Mahasiswa with ID = %s was updated", params["id"])
}

func deleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Transaction declaration
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check if mahasiswa exists in database
	var mahasiswa Mahasiswa
	err = tx.QueryRow("SELECT id FROM mahasiswa WHERE id = ?", params["id"]).Scan(&mahasiswa.ID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			http.Error(w, "Mahasiswa not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the corresponding rows in mahasiswa_hobi table
	hobiStmt, err := tx.Prepare("DELETE FROM mahasiswa_hobi WHERE id_mahasiswa =?")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = hobiStmt.Exec(params["id"])
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the mahasiswa

	mahasiswaStmt, err := tx.Prepare("DELETE FROM mahasiswa WHERE id =?")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = mahasiswaStmt.Exec(params["id"])
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Mahasiswa with ID = %s and associated hobbies were deleted", params["id"])
}
