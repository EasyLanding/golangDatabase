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

func CreateUserTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            age INTEGER,
            email TEXT,
            username TEXT
        )`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	fmt.Println("User table created successfully")
	return nil
}

func InsertUser(user User, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

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

func SelectUser(userID int, db *sql.DB) (User, error) {
	if db == nil {
		return User{}, fmt.Errorf("db is nil")
	}

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

func UpdateUser(user User, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

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

func DeleteUser(userID int, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

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

func SelectAllUsers(db *sql.DB) ([]User, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

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
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка открытия базы данных:", err)
		return
	}

	defer db.Close()

	fmt.Println(HelloWorld())

	// Создание таблицы пользователей
	err = CreateUserTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	userID := 2
	user, err := SelectUser(userID, db)
	if err != nil {
		fmt.Println("Ошибка при выборе пользователя:", err)
		return
	}
	fmt.Println("Выбранный пользователь:")
	fmt.Println("ID:", user.ID)
	fmt.Println("Username:", user.Username)
	fmt.Println("Email:", user.Email)

	//Обновление юзера
	user.Username = "jane_doe"
	user.Email = "jane.doe@example.com"
	err = UpdateUser(user, db)
	if err != nil {
		fmt.Println("Ошибка при вставке пользователя:", err)
		return
	}

	//Удаляем юзера
	err = DeleteUser(userID, db)
	if err != nil {
		fmt.Println("Ошибка при вставке пользователя:", err)
		return
	}

	//Получаем всех юзхеров из бд
	allUsers, err := SelectAllUsers(db)
	if err != nil {
		fmt.Println("Ошибка при получении всех пользователей:", err)
	} else {
		fmt.Println("Все пользователи:")
		for _, user := range allUsers {
			fmt.Println(user)
		}
	}
}
