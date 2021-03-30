package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"path"

	_ "github.com/lib/pq"
)

type User struct {
	Id            int
	Date          string
	VerifyMode    string
	Name          string
	PicturesCount int
	Blob          string
}

func main() {
	http.HandleFunc("/", ShowUsers)
	http.ListenAndServe(":8080", nil)
}

func ShowUsers(w http.ResponseWriter, r *http.Request) {
	users := getUsers()
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUsers() []User {
	connStr := "user=postgres host=127.0.0.1 password=XGalHeg7 dbname=golangdb sslmode=disable port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from users;")
	if err != nil {
		panic(err)
	}
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Date, &user.VerifyMode, &user.Name, &user.PicturesCount, &user.Blob)
		if err != nil {
			fmt.Println(err)
			continue
		}
		user.Blob = base64.StdEncoding.EncodeToString([]byte(user.Blob))
		users = append(users, user)
	}
	return users
}
