package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jobhun-test/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MahasiswaHandler struct {
	repo repository.MahasiswaRepository
}

func NewMahasiswaHandler(db *sql.DB) *MahasiswaHandler {
	return &MahasiswaHandler{
		repo: repository.NewMahasiswaRepository(db),
	}
}

func (h *MahasiswaHandler) GetMahasiswaById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	mahasiswa, err := h.repo.GetMahasiswaById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving mahasiswa data"})
		return
	}
	if mahasiswa == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Mahasiswa not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mahasiswa)
}
func (h *MahasiswaHandler) GetAllMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMahasiswa, err := h.repo.GetAllMahasiswa()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving mahasiswa data"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if len(allMahasiswa) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No mahasiswa data found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allMahasiswa)
}

func (h *MahasiswaHandler) InsertMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mahasiswa := repository.Mahasiswa{}
	err := json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	mahasiswaInserted, err := h.repo.InsertMahasiswa(&mahasiswa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":           "new mahasiswa was created",
		"insertedMahasiswa": mahasiswaInserted,
	})
	fmt.Fprintf(w, "New mahasiswa was created")

	if err != nil {
		http.Error(w, "failed to insert mahasiswa", http.StatusInternalServerError)
		return
	}

}

func (h *MahasiswaHandler) UpdateMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idMahasiswa, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid mahasiswa id", http.StatusBadRequest)
		return
	}

	mahasiswa := repository.Mahasiswa{}
	err = json.NewDecoder(r.Body).Decode(&mahasiswa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedMahasiswa, err := h.repo.UpdateMahasiswa(&mahasiswa, idMahasiswa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMahasiswa)
}

func (h *MahasiswaHandler) DeleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	fmt.Printf("The type of variable is : %T\n", id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}
	err = h.repo.DeleteMahasiswa(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Mahasiswa not found", http.StatusNotFound)
			return
		}
		fmt.Println("here")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mahasiswa with ID = %d and associated hobbies were deleted", id)
}
