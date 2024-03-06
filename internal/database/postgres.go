package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/say8hi/go-api-test/internal/models"
)

var db *sql.DB

func Init() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			fmt.Println("Соединение с базой данных успешно установлено")
			return
		}
		fmt.Printf("Не удалось подключиться к базе данных: %v. Повторная попытка через 1 секунду\n", err)
		time.Sleep(1 * time.Second)
	}
}

func CreateTables() {
	tables := [4]string{
		`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT NOT NULL,
            full_name TEXT,
            password_hash TEXT NOT NULL
        );
    `,
		`
        CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255),
            description TEXT
        );
    `,
		`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255),
            description TEXT,
            price NUMERIC(10,2)
        );
    `,
		`
        CREATE TABLE IF NOT EXISTS product_category (
            product_id INT,
            category_id INT,
            FOREIGN KEY (product_id) REFERENCES products(id),
            FOREIGN KEY (category_id) REFERENCES categories(id),
            PRIMARY KEY (product_id, category_id)
        );
    `,
	}

	for _, sql := range tables {
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CreateUser(request_user models.CreateUserRequest) (models.UserInDatabase, error) {
	var user models.UserInDatabase
	err := db.QueryRow("INSERT INTO users (username, full_name, password_hash) VALUES($1, $2, $3) RETURNING *",
		request_user.Username, request_user.FullName, request_user.Password).Scan(&user.ID, &user.Username, &user.FullName, &user.PasswordHash)
	if err != nil {
		return models.UserInDatabase{}, err
	}
	return user, nil
}
