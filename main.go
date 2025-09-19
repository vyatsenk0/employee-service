package main

import (
	"context"
	"employee-service/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"        // для REST API
	"github.com/jackc/pgx/v5/pgxpool" // для подключения к Postgres
	"github.com/joho/godotenv"        // для загрузки .env
)

var db *pgxpool.Pool // для пула соединений к postgresql

// загрузить переменную DATABASE_URL из .env для доступа через os.Getenv()
func init() {
	_ = godotenv.Load()
}

// создать новый роутер Gin и зарегистрировать POST эндпоинт
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/employees", createEmployee)
	return r
}

func main() {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	log.Println("Connecting to:", dsn)

	var err error

	db, err = pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Unable to connect to database : %v\n", err)
	}
	defer db.Close()

	initDB() // автоматически создать таблицу, если нет

	//r := gin.Default() // создать роутер Gin
	//r.POST("/employees", createEmployee)
	r := setupRouter()
	r.Run(":8080")
}

// создать таблицу, если еще нет
func initDB() {
	query := `CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		phone VARCHAR(20),
		city VARCHAR(50)
	);`
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
}

// POST /employees
func createEmployee(c *gin.Context) {
	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil { // парсинг json из запроса в структуру emp
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверить, что все поля заполнены
	if emp.Name == "" || emp.Phone == "" || emp.City == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields must be filled"})
		return
	}

	// вставить запись в таблицу и вернуть новый id
	query := `INSERT INTO employees (name, phone, city) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(context.Background(), query, emp.Name, emp.Phone, emp.City).Scan(&emp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// вернуть json с созданным обьектом и id
	c.JSON(http.StatusOK, emp)
}
