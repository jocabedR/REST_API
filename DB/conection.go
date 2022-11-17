package main

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectPostgres() *sql.DB {
	Hostname := "localhost"
	Port := 5432
	Username := "joca"
	Password := "joca"
	Database := "restapi"
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%ssslmode=disable", Hostname, Port, Username, Password, Database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
		return nil
	}
	return db
}

func DeleteUser(ID int) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}
	defer db.Close()
	t := FindUserID(ID)
	if t.ID == 0 {
		log.Println("User", ID, "does not exist.")
		return false
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE ID = $1")
	if err != nil {
		log.Println("DeleteUser:", err)
		return false
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		log.Println("DeleteUser:", err)
		return false
	}
	return true
}

func ListAllUsers() []User {
	db := ConnectPostgres()
	if db == nil {
		fmt.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return []User{}
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users \n")
	if err != nil {
		log.Println(err)
		return []User{}
	}

	all := []User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		temp := User{c1, c2, c3, c4, c5, c6}
		all = append(all, temp)
	}
	log.Println("All:", all)
	return all
}

func IsUserValid(u User) bool {
	db := ConnectPostgres()
	if db == nil {
		fmt.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE Username = $1 \n",
		u.Username)
	if err != nil {
		log.Println(err)
		return false
	}
	temp := User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = User{c1, c2, c3, c4, c5, c6}
	}
	if u.Username == temp.Username && u.Password == temp.Password {
		return true
	}

	return false
}
