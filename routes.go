package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, GROUP_CONCAT(h.nama_hobi SEPARATOR ';') FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id WHERE m.id =? GROUP BY m.id", params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving mahasiswa data"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		panic(err.Error())
	}
	defer result.Close()

	var mahasiswa Mahasiswa
	var hobi string
	for result.Next() {
		err := result.Scan(&mahasiswa.ID,
			&mahasiswa.Nama,
			&mahasiswa.Usia,
			&mahasiswa.Gender,
			&mahasiswa.Tanggal_Registrasi,
			&mahasiswa.ID_Jurusan,
			&hobi)
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
		// Split comma-separated string of hobbies into a slice of strings
		hobiArr := strings.Split(hobi, ";")
		mahasiswa.Hobi = hobiArr

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mahasiswa)
	}

}

func getAllMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var allMahasiswa []Mahasiswa

	result, err := db.Query("SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, GROUP_CONCAT(h.nama_hobi SEPARATOR ';') FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id GROUP BY m.id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving mahasiswa data"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		panic(err.Error())
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

		// Split comma-separated string of hobbies into a slice of strings
		hobiArr := strings.Split(hobi, ";")
		mahasiswa.Hobi = hobiArr

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

	var mahasiswa Mahasiswa
	// err = json.Unmarshal(body, &mahasiswa)
	err = json.NewDecoder(r.Body).Decode(&mahasiswa)
	fmt.Print(mahasiswa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	vars := mux.Vars(r)
	idMahasiswa, err := strconv.Atoi(vars["id"])
	println(idMahasiswa)
	if err != nil {
		http.Error(w, "invalid mahasiswa id", http.StatusBadRequest)
		return
	}

	stmt, err := tx.Prepare("UPDATE mahasiswa SET nama = ?, usia = ?, gender = ?, tanggal_registrasi = ?, id_jurusan = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "failed to prepare update statement", http.StatusInternalServerError)
		return
	}

	var mahasiswa Mahasiswa
	err = json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, mahasiswa.ID_Jurusan, idMahasiswa)
	if err != nil {
		http.Error(w, "failed to execute update statement ", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM mahasiswa_hobi WHERE id_mahasiswa = ?", idMahasiswa)
	if err != nil {
		http.Error(w, "failed to delete existing mahasiswa_hobi data", http.StatusInternalServerError)
		return
	}

	for _, h := range mahasiswa.Hobi {
		stmt, err = tx.Prepare("INSERT INTO mahasiswa_hobi(id_mahasiswa, id_hobi) VALUES (?, ?)")
		if err != nil {
			http.Error(w, "failed to prepare insert statement for mahasiswa_hobi table", http.StatusInternalServerError)
			panic(err.Error())
		}

		_, err = stmt.Exec(idMahasiswa, h)
		if err != nil {
			http.Error(w, "failed to execute query", http.StatusInternalServerError)
			panic(err.Error())

		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "mahasiswa data was updated"})
	fmt.Fprintf(w, "Mahasiswa data was updated")
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
