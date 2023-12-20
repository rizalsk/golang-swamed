package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver MYSql
	_ "github.com/mattn/go-sqlite3"    // Driver SQLite
)

const (
	dbconf   = "golangcrud.db"
	dbdriver = "sqlite3" //"mysql"
	note     = `
			==== Aplikasi basic CRUD GoLang ====
Pilih:
0. Berhenti
1. Create
2. Read
3. Update
4. Delete
5. lanjutkan`
)

type user struct {
	id    uint
	name  string
	email string
}

func main() {
	defer catch()
	option()
}

func create() {
	var (
		name  = "user"
		email = "user@mail.com"
	)
	db, err := sql.Open("sqlite3", "nama_db.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users VALUES (?, ?)", name, email)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("insert success!")
}

func read() {
	db, err := sql.Open("sqlite3", "nama_db.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string

		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", id, name, email)
	}
}

func update() {
	var (
		id    uint   = 1
		name  string = "userUpdate"
		email string = "userUpdate@mail.com"
	)

	db, err := sql.Open("sqlite3", "nama_db.db")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users set name = ?, email = ? where id = ?", name, email, id)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func delete() {
	id := 1
	db, err := sql.Open("sqlite3", "nama_db.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getByID(id uint, db *sql.DB) (user, error) {
	sqlStatement := `SELECT * FROM users WHERE id=$1;`

	var existingUser user

	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&existingUser.id, &existingUser.name, &existingUser.email); err {
	case sql.ErrNoRows:
		return existingUser, errors.New("Data tidak ditemukan")
	case nil:
		return existingUser, nil
	default:
		return existingUser, nil
	}
}

func catch() {
	if r := recover(); r != nil {
		fmt.Println("Error occured", r)
	} else {
		fmt.Println("Application running perfectly")
	}
}
