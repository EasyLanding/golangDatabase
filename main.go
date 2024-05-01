package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func HelloWorld() string {
	return "Hello world!"
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func CreateUserTable() error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt :=
		`CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  age INTEGER
 )`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	fmt.Println("User table created successfully")
	return nil
}

func InsertUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt :=
		`INSERT INTO users (name, age) VALUES (?, ?)`

	_, err = db.Exec(sqlStmt, user.Name, user.Age)
	if err != nil {
		return err
	}

	fmt.Println("User inserted successfully")
	return nil
}

func SelectUser(id int) (User, error) {
	var user User

	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return user, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id)
	err = row.Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt :=
		`UPDATE users SET name = ?, age = ? WHERE id = ?`

	_, err = db.Exec(sqlStmt, user.Name, user.Age, user.ID)
	if err != nil {
		return err
	}

	fmt.Println("User updated successfully")
	return nil
}

func DeleteUser(id int) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt :=
		`DELETE FROM users WHERE id = ?`

	_, err = db.Exec(sqlStmt, id)
	if err != nil {
		return err
	}

	fmt.Println("User deleted successfully")
	return nil
}

func SelectAllUsers() ([]User, error) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func main() {
	fmt.Println(HelloWorld())

	allUsers, err := SelectAllUsers()
	if err != nil {
		fmt.Println("Ошибка при получении всех пользователей:", err)
	} else {
		fmt.Println("Все пользователи:")
		for _, user := range allUsers {
			fmt.Println(user)
		}
	}
}
