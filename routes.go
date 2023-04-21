package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, h.nama_hobi FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id WHERE m.id =?", params["id"])
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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mahasiswa)
	}

}

func getAllMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var allMahasiswa []Mahasiswa

	result, err := db.Query("SELECT m.id, m.nama, m.usia, m.gender, m.tanggal_registrasi, j.nama_jurusan, h.nama_hobi FROM mahasiswa m INNER JOIN jurusan j ON m.id_jurusan = j.id LEFT JOIN mahasiswa_hobi hm ON m.id = hm.id_mahasiswa LEFT JOIN hobi h ON hm.id_hobi = h.id")
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

	stmt, err := db.Prepare("INSERT INTO mahasiswa(nama, usia, gender, tanggal_registrasi, id_jurusan) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to prepare insert statement"})
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to read request body"})
		panic(err.Error())
	}

	var mahasiswa Mahasiswa
	err = json.Unmarshal(body, &mahasiswa)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse request body"})
		panic(err.Error())
	}

	result, err := stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, mahasiswa.ID_Jurusan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to execute insert statement"})
		panic(err.Error())
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get last inserted ID"})
		panic(err.Error())
	}
	for _, h := range mahasiswa.Hobi {
		stmt, err = db.Prepare("INSERT INTO mahasiswa_hobi(id_mahasiswa, id_hobi) VALUES (?, ?)")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to execute insert statement for mahasiswa_hobi table"})
			panic(err.Error())
		}

		_, err = stmt.Exec(lastID, h)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to execute query"})
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

	// Delete the corresponding rows in mahasiswa_hobi table
	hobiStmt, err := db.Prepare("DELETE FROM mahasiswa_hobi WHERE id_mahasiswa =?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = hobiStmt.Exec(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the mahasiswa

	mahasiswaStmt, err := db.Prepare("DELETE FROM mahasiswa WHERE id =?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = mahasiswaStmt.Exec(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Mahasiswa with ID = %s and associated hobbies were deleted", params["id"])
}
