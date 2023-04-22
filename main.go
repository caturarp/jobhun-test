package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Mahasiswa struct {
	ID                 int    `json:"id"`
	Nama               string `json:"nama"`
	Usia               int    `json:"usia"`
	Gender             string `json:"gender"`
	Tanggal_Registrasi string `json:"tanggal_registrasi"`
	ID_Jurusan         string `json:"id_jurusan"`
	Hobi               string `json:"hobi"`
}

// type Hobi struct {
// 	ID   int    `json:"id"`
// 	Nama string `json:"nama"`
// }

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@/jobhun")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to MySQL!")
	router := mux.NewRouter()

	router.HandleFunc("/students", insertMahasiswa).Methods("POST")     // add new mahasiswa
	router.HandleFunc("/students/{id}", updateMahasiswa).Methods("PUT") // edit detail mahasiswa
	router.HandleFunc("/students", getAllMahasiswa).Methods("GET")      // mahasiswa +jurusan+hobi
	router.HandleFunc("/students/{id}", getMahasiswa).Methods("GET")
	router.HandleFunc("/students/{id}", deleteMahasiswa).Methods("DELETE")
	http.ListenAndServe("127.0.0.1:8000", router)

}
