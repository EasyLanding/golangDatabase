package main

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

func HelloWorld() string {
	return "Hello world!"
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
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
            age INTEGER,
            email TEXT,
            username TEXT
        )`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	fmt.Println("User table created successfully")
	return nil
}

func PrepareQuery(operation string, table string, user User) (string, []interface{}, error) {
	var query string
	var args []interface{}

	switch operation {
	case "insert":
		query, args, _ = squirrel.Insert(table).Columns("username", "email").Values(user.Username, user.Email).ToSql()
	case "select":
		query, args, _ = squirrel.Select("id", "username", "email").From(table).Where(squirrel.Eq{"id": user.ID}).ToSql()
	case "update":
		query, args, _ = squirrel.Update(table).Set("username", user.Username).Set("email", user.Email).Where(squirrel.Eq{"id": user.ID}).ToSql()
	case "delete":
		query, args, _ = squirrel.Delete(table).Where(squirrel.Eq{"id": user.ID}).ToSql()
	default:
		return "", nil, fmt.Errorf("invalid operation")
	}

	return query, args, nil
}

func InsertUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("insert", "users", user)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func SelectUser(userID int) (User, error) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	var user User
	query, args, err := PrepareQuery("select", "users", User{ID: userID})
	if err != nil {
		return User{}, err
	}

	row := db.QueryRow(query, args...)
	err = row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UpdateUser(user User) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("update", "users", user)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(userID int) error {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query, args, err := PrepareQuery("delete", "users", User{ID: userID})
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func SelectAllUsers() ([]User, error) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func main() {
	fmt.Println(HelloWorld())

	var err error

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
