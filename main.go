package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"jobhun-test/handler"

	_ "jobhun-test/docs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// @title			Jobhun API
// @version		1
// @description	This is the Jobhun API documentation.
// @host			localhost:8000
// @BasePath		/

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

	mahasiswaHandler := handler.NewMahasiswaHandler(db)

	router := mux.NewRouter()

	router.HandleFunc("/students", mahasiswaHandler.InsertMahasiswa).Methods("POST")     // add new mahasiswa
	router.HandleFunc("/students/{id}", mahasiswaHandler.UpdateMahasiswa).Methods("PUT") // edit detail mahasiswa
	router.HandleFunc("/students", mahasiswaHandler.GetAllMahasiswa).Methods("GET")      // mahasiswa +jurusan+hobi
	router.HandleFunc("/students/{id}", mahasiswaHandler.GetMahasiswaById).Methods("GET")
	router.HandleFunc("/students/{id}", mahasiswaHandler.DeleteMahasiswa).Methods("DELETE")
	// router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Starting server on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", router))
}
