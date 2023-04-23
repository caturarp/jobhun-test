# Jobhun Backend Internship Test

This is the backend internship test for Jobhun. The goal of this test is to create a simple RESTful API for managing data on Mahasiswa (students) in a university.

## Getting Started

### Requirements

- Go 1.16 or higher
- MySQL

### Install dependencies:

```bash
git clone https://https://github.com/caturarp/jobhun-test.git
cd jobhun-test
go mod download
```
### Build and Run
```bash
go build
./jobhun-test.exe
```
## API Documentation

The following API endpoints are available:

### GET /students/{id}

Get a single mahasiswa by ID.

```bash 
curl http://localhost:8080/mahasiswa/1
```
### GET /students/
Get all mahasiswa.
```bash
curl http://localhost:8080/mahasiswa
```

### PUT /students/{id}
Edit info mahasiswa
```bash
curl -X PUT -H "Content-Type: application/json" -d '{"Nama": "John Doe Jr.", "Usia": 21, "Gender": 1, "Tanggal_Registrasi": "2022-04-23", "ID_Jurusan": 3, "Hobi": ["4", "3"]}' http://localhost:8080/students/1
```

### POST /students
Insert a new mahasiswa.
```bash
curl -X POST -H "Content-Type: application/json" -d '{"Nama": "Jehn Dal", "Usia": 20, "Gender": "0", "Tanggal_Registrasi": "2022-04-23", "ID_Jurusan": 1, "Hobi" : ["2","4"]}' http://localhost:8080/students
```

### DELETE /students/{id}
Delete 1 mahasiswa
```bash
curl -X DELETE http://localhost:8000/students/1
```