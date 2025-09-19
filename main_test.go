package main // для теста кода из main.go

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert" // для тестов
)

var testDB *pgxpool.Pool

// Подключение к тестовой БД и очистка таблицы
func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL not set")
	}

	var err error
	testDB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// создать таблицу, если нет
	_, err = testDB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			phone VARCHAR(20),
			city VARCHAR(50)
		);
	`)
	if err != nil {
		panic(err)
	}

	// очистить таблицу перед запуском тестов
	_, err = testDB.Exec(context.Background(), `TRUNCATE employees RESTART IDENTITY`)

	if err != nil {
		panic(err)
	}

	// тестовую БД в глобальную переменную из main.go
	db = testDB

	code := m.Run()
	os.Exit(code)
}

// проверка на пустой json request
func TestCreateEmployee_BadRequest(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/employees", strings.NewReader(`{}`)) // пустой JSON
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// проверка POST /employees
func TestCreateEmployee(t *testing.T) {
	router := setupRouter() // роутер gin.Engine

	w := httptest.NewRecorder()
	reqBody := `{"name":"John","phone":"123456","city":"Almaty"}`
	req, _ := http.NewRequest("POST", "/employees", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	t.Logf("Response body: %s", w.Body.String()) // вывод ответа в json

	assert.Equal(t, 200, w.Code)                // проверка статуса
	assert.Contains(t, w.Body.String(), "John") // проверка возврата

	// проверить, что запись появилась в БД
	var count int
	err := testDB.QueryRow(context.Background(), "SELECT COUNT(*) FROM employees WHERE name=$1", "John").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
