package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
)

func TestHelloWorldFunc(t *testing.T) {
	var buf bytes.Buffer

	fmt.Fprint(&buf, HelloWorld())

	output := buf.String()
	expected := "Hello world!"

	if output != expected {
		t.Errorf("Unexpected output: %s", output)
	}
}

func TestPrepareQuery(t *testing.T) {
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при открытии базы данных: %s", err)
	}
	defer db.Close()

	query, args, err := PrepareQuery("insert", "users", user)
	if err != nil {
		t.Errorf("Ошибка при подготовке запроса insert: %s", err)
	}

	if query == "" {
		t.Error("Пустой запрос insert")
	}

	if len(args) != 2 {
		t.Error("Неверное количество аргументов для insert")
	}

	// Другие тесты для остальных операций могут быть добавлены аналогичным образом
}

func TestCreateUserTable(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при открытии базы данных: %s", err)
	}
	defer db.Close()

	err = CreateUserTable(db)
	if err != nil {
		t.Errorf("Ошибка при создании таблицы пользователей: %s", err)
	}

	// Дополнительные проверки, что таблица была создана успешно
}

func TestInsertUser(t *testing.T) {
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при открытии базы данных: %s", err)
	}
	defer db.Close()

	err = CreateUserTable(db)
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы пользователей: %s", err)
	}

	err = InsertUser(user, db)
	if err != nil {
		t.Errorf("Ошибка при вставке пользователя: %s", err)
	}

	// Дополнительные проверки, что пользователь был успешно добавлен
}

func TestSelectUser(t *testing.T) {
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при открытии базы данных: %s", err)
	}
	defer db.Close()

	err = CreateUserTable(db)
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы пользователей: %s", err)
	}

	err = InsertUser(user, db)
	if err != nil {
		t.Fatalf("Ошибка при вставке пользователя: %s", err)
	}

	selectedUser, err := SelectUser(user.ID, db)
	if err != nil {
		t.Errorf("Ошибка при выборке пользователя: %s", err)
	}

	if selectedUser.Username != user.Username || selectedUser.Email != user.Email {
		t.Errorf("Данные пользователя не совпадают")
	}
}

func TestUpdateUser(t *testing.T) {
	// Создаем тестового пользователя
	testUser := User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при открытии базы данных: %s", err)
	}
	defer db.Close()

	// Создаем таблицу пользователей для тестирования
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, email TEXT)")
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы: %s", err)
	}

	// Вставляем тестового пользователя
	_, err = db.Exec("INSERT INTO users (username, email) VALUES ('testuser', 'testuser@example.com')")
	if err != nil {
		t.Fatalf("Ошибка при вставке тестового пользователя: %s", err)
	}

	err = UpdateUser(testUser, db)

	if err != nil {
		t.Errorf("Ошибка при обновлении пользователя: %s", err)
	}

	// Проверяем, что данные пользователя были успешно обновлены
	var username string
	var email string
	err = db.QueryRow("SELECT username, email FROM users WHERE id = ?", testUser.ID).Scan(&username, &email)
	if err != nil {
		t.Fatalf("Ошибка при получении данных пользователя: %s", err)
	}

	if username != testUser.Username || email != testUser.Email {
		t.Errorf("Данные пользователя не были обновлены")
	}
}

func TestDeleteUser(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при открытии базы данных: %s", err)
	}
	defer db.Close()

	// Создаем таблицу пользователей для тестирования
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, email TEXT)")
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы: %s", err)
	}

	// Вставляем тестового пользователя
	_, err = db.Exec("INSERT INTO users (username, email) VALUES ('testuser', 'testuser@example.com')")
	if err != nil {
		t.Fatalf("Ошибка при вставке тестового пользователя: %s", err)
	}

	// Вызываем тестируемую функцию
	err = DeleteUser(1, db) // Передаем ID пользователя для удаления

	// Проверяем ошибки
	if err != nil {
		t.Errorf("Ошибка при удалении пользователя: %s", err)
	}

	// Проверяем, что пользователь был удален
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Ошибка при подсчете пользователей: %s", err)
	}
	if count != 0 {
		t.Errorf("Пользователь не был удален")
	}
}

func TestSelectAllUsers(t *testing.T) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		fmt.Println("Ошибка открытия базы данных:", err)
		return
	}

	defer db.Close()

	// Вызываем тестируемую функцию
	users, err := SelectAllUsers(db)

	// Проверяем ошибки
	if err != nil {
		t.Errorf("Ошибка при выборе всех пользователей: %s", err)
	}

	// Проверяем количество пользователей
	if len(users) != 2 {
		t.Errorf("Ожидается 2 пользователя, получено: %d", len(users))
	}

	// Проверяем данные первого пользователя
	if users[0].ID != 1 || users[0].Username != "jane_doe" || users[0].Email != "jane.doe@example.com" {
		t.Errorf("Некорректные данные первого пользователя")
	}

	// Проверяем данные второго пользователя
	if users[1].ID != 3 || users[1].Username != "alice123" || users[1].Email != "alice@example.com" {
		t.Errorf("Некорректные данные второго пользователя")
	}
}
