package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type User struct {
	ID           int
	Name         string
	RegTimeStamp string
	AccessFlags  int
}

var database *sql.DB

func GetUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var Users = []User{}
	rows, err := database.Query("select * from new_schema.users")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := User{}
		err := rows.Scan(&p.ID, &p.Name, &p.RegTimeStamp, &p.AccessFlags)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Users = append(Users, p)
	}

	b, err := json.Marshal(Users)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	p := User{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Decoded to struct ", p)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = database.Exec("insert into new_schema.users values (?, ?, ?, ?)",
		&p.ID,
		&p.Name,
		&p.RegTimeStamp,
		&p.AccessFlags)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusLocked)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func main() {

	db, err := sql.Open("mysql", "root:qwerty123@/new_schema")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users", PostUser).Methods(http.MethodPost)
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	_ = http.ListenAndServe(":8181", nil)
}
