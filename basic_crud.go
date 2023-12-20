package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver MYSql
	_ "github.com/mattn/go-sqlite3"    // Driver SQLite
)

const (
	dbconf = "golangcrud.db"
	// dbconf   = "root:@tcp(localhost:3306)/golangcrud"
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

func option() {
	var pilihan int
	fmt.Println(note)
	fmt.Print("Pilihan:")
	fmt.Scan(&pilihan)
	switch pilihan {
	case 1:
		create()
	case 2:
		read()
	case 3:
		update()
	case 4:
		delete()
	case 5:
		option()
	default:
		fmt.Println("program berhenti")
		os.Exit(1)
	}
}

func create() {
	var input user
	db, err := sql.Open(dbdriver, dbconf)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Print("==== input user baru === \nname: ")
	fmt.Scan(&input.name)
	fmt.Print("email: ")
	fmt.Scan(&input.email)

	sqlStatement := `
		INSERT INTO users
		(name, email)
		VALUES($1,$2);`

	result, err := db.Exec(sqlStatement, input.name, input.email)
	if err != nil {
		panic(err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Data berhasil ditambahkan!")
	option()
}

func read() {
	db, err := sql.Open(dbdriver, dbconf)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var (
		users []user
		exist = false
	)

	for rows.Next() {
		var each = user{}
		if err := rows.Scan(&each.id, &each.name, &each.email); err != nil {
			log.Fatal(err)
		}
		users = append(users, each)
		exist = true
	}

	if exist {
		fmt.Println("\n\t===Data Users===\t")
		for _, each := range users {
			fmt.Printf("ID: %d, Name: %s, Email: %s\n", each.id, each.name, each.email)
		}
	} else {
		fmt.Println("\nTidak ada record ditemukan")
	}

	option()
}

func update() {
	var (
		id           uint
		existingUser user
		userUpdate   user
	)

	db, err := sql.Open(dbdriver, dbconf)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Print("==== data update user === \nuser dengan id berapa yang akan diupdate: ")
	fmt.Scan(&id)

	// getUserByID
	existingUser, err = getByID(id, db)
	if err != nil {
		// panic(err.Error())
		fmt.Println(err.Error())

	} else {
		if existingUser.id == 0 {
			fmt.Println("Record tidak ditemukan")
		} else {
			fmt.Printf("ID: %d, Name: %s, Email: %s\n", existingUser.id, existingUser.name, existingUser.email)

			fmt.Print("update name: ")
			fmt.Scan(&userUpdate.name)
			fmt.Print("update email: ")
			fmt.Scan(&userUpdate.email)

			// jika data Name kosong
			if userUpdate.name == "" {
				userUpdate.name = existingUser.name
			}

			// jika data Email kosong
			if userUpdate.email == "" {
				userUpdate.email = existingUser.email
			}

			res, err := db.Exec("UPDATE users set name = ?, email = ? where id = ?", userUpdate.name, userUpdate.email, id)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			count, err := res.RowsAffected()
			if count >= 1 {
				fmt.Println("Data berhasil diperbarui!")
			} else {
				fmt.Println("Data gagal diperbarui! error: ", err.Error())
			}
		}
	}
	option()
}

func delete() {
	var (
		id           uint
		existingUser user
	)

	db, err := sql.Open(dbdriver, dbconf)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Print("==== pilih id user yang akan dihapus === \nID User: ")
	fmt.Scan(&id)

	existingUser, err = getByID(id, db)
	if err != nil {
		// panic(err.Error())
		fmt.Println(err.Error())
		return
	}

	fmt.Println("data yang akan dihapus: ", existingUser)

	sqlStatement := `
	DELETE FROM users
	WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Data berhasil dihapus!")
	option()
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
