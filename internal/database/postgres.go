package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/say8hi/go-api-test/internal/models"
)

var db *sql.DB

var ErrCreatingProduct = errors.New("error creating product")
var ErrRollback = errors.New("error deleting product")
var ErrCategoryDoesntExists = errors.New("error category from categories field doesn't exists")

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
			fmt.Println("Successfully connected to database")
			return
		}
		fmt.Printf("Couldn't connect to databse: %v. Retry in 1 second.\n", err)
		time.Sleep(1 * time.Second)
	}
}

func CloseConnection() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Соединение с базой данных успешно закрыто")
	}
}

func CreateTables() {
	tables := [4]string{
		`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT NOT NULL UNIQUE,
            full_name TEXT,
            password_hash TEXT NOT NULL
        );
    `,
		`
        CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) UNIQUE,
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
            FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
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

// Table Users
func CreateUser(request_user models.CreateUserRequest) (models.UserInDatabase, error) {
	var user models.UserInDatabase
	err := db.QueryRow("INSERT INTO users (username, full_name, password_hash) VALUES($1, $2, $3) RETURNING *",
		request_user.Username, request_user.FullName, request_user.Password).Scan(&user.ID, &user.Username, &user.FullName, &user.PasswordHash)

	if err != nil {
		return models.UserInDatabase{}, err
	}

	return user, nil
}

func GetUserByUsername(username string) (models.UserInDatabase, error) {
	var user models.UserInDatabase

	err := db.QueryRow("SELECT * FROM users WHERE username=$1",
		username).Scan(&user.ID, &user.Username, &user.FullName, &user.PasswordHash)

	if err != nil {
		return models.UserInDatabase{}, err
	}

	return user, nil
}

func GetUserByPasswordHash(password string) (models.UserInDatabase, error) {
	var user models.UserInDatabase

	err := db.QueryRow("SELECT * FROM users WHERE password_hash=$1",
		password).Scan(&user.ID, &user.Username, &user.FullName, &user.PasswordHash)

	if err != nil {
		return models.UserInDatabase{}, err
	}

	return user, nil
}

// Table Categories
func CreateCategory(createCategory models.CreateCategoryRequest) (models.Category, error) {
	var category models.Category

	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description`
	err := db.QueryRow(query, &createCategory.Name, &createCategory.Description).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return models.Category{}, fmt.Errorf("error creating category: %v", err)
	}

	return category, nil
}

func GetCategoryByID(category_id int) (models.Category, error) {
	var category models.Category

	query := `SELECT * FROM categories WHERE id=$1`
	err := db.QueryRow(query, category_id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func GetAllCategories() ([]models.Category, error) {
	categories := []models.Category{}
	query := `SELECT * FROM categories ORDER BY id`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error creating category: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

func UpdateCategory(categoryID int, updateReq models.CategoryUpdateRequest) error {
	var setParts []string
	var args []interface{}
	var argIndex int = 1

	if updateReq.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *updateReq.Name)
		argIndex++
	}
	if updateReq.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *updateReq.Description)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClause := strings.Join(setParts, ", ")
	queryString := fmt.Sprintf("UPDATE categories SET %s WHERE id = $%d", setClause, argIndex)
	args = append(args, categoryID)

	_, err := db.Exec(queryString, args...)
	if err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}

	return nil
}

func DeleteCategory(categoryID int) error {
	queryString := "DELETE FROM categories WHERE id = $1"
	_, err := db.Exec(queryString, categoryID)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}

	return nil
}

// Table Products
func CreateProduct(productRequest models.CreateProductRequest) (models.Product, error) {
    var product models.Product

    tx, err := db.Begin()
    if err != nil {
        return models.Product{}, err
    }

    productQuery := `INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id`
    err = tx.QueryRow(productQuery, productRequest.Name, productRequest.Description, productRequest.Price).Scan(&product.ID)
    if err != nil {
        tx.Rollback()
        return models.Product{}, ErrCreatingProduct
    }

    var categories []models.Category
    for _, categoryName := range productRequest.Categories {
        var category models.Category
        categoryQuery := `SELECT id, name, description FROM categories WHERE name = $1`
        err := tx.QueryRow(categoryQuery, categoryName).Scan(&category.ID, &category.Name, &category.Description)
        if err != nil {
            tx.Rollback()
            return models.Product{}, ErrCategoryDoesntExists
        }

        _, err = tx.Exec(`INSERT INTO product_category (product_id, category_id) VALUES ($1, $2)`, product.ID, category.ID)
        if err != nil {
            tx.Rollback()
            return models.Product{}, ErrCreatingProduct
        }
        categories = append(categories, category)
    }

    if err := tx.Commit(); err != nil {
        return models.Product{}, err
    }

    product.Name = productRequest.Name
    product.Description = productRequest.Description
    product.Price = productRequest.Price
    product.Categories = categories

    return product, nil
}

func GetProduct(productId int) (models.Product, error) {
	var product models.Product

	productQuery := `SELECT id, name, description, price FROM products WHERE id = $1`
	err := db.QueryRow(productQuery, productId).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		return models.Product{}, err
	}

	categoriesQuery := `
SELECT c.id, c.name, c.description
FROM categories c
INNER JOIN product_category pc ON c.id = pc.category_id
WHERE pc.product_id = $1 ORDER BY c.id ASC
`
	rows, err := db.Query(categoriesQuery, productId)
	if err == sql.ErrNoRows {
		return product, nil
	} else if err != nil {
    return product, fmt.Errorf("error fetching categories for product: %v", err)
  }
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return models.Product{}, fmt.Errorf("error scanning category: %v", err)
		}
		product.Categories = append(product.Categories, category)
	}
	if err := rows.Err(); err != nil {
		return models.Product{}, fmt.Errorf("error iterating categories: %v", err)
	}

	return product, nil
}

func GetProductsByCategory(categoryID int) ([]models.Product, error) {
	query := `
SELECT p.id, p.name, p.description, p.price, c.id, c.name, c.description
FROM products p
JOIN product_category pc ON p.id = pc.product_id
JOIN categories c ON pc.category_id = c.id
WHERE pc.category_id = $1
ORDER BY p.id, c.id
`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error querying products by category: %v", err)
	}
	defer rows.Close()

	productsMap := make(map[int]*models.Product)
	for rows.Next() {
		var (
			prodID, catID   int
			prodName, prodDescription string
			prodPrice       float64
			catName, catDescription string
		)
		if err := rows.Scan(&prodID, &prodName, &prodDescription, &prodPrice, &catID, &catName, &catDescription); err != nil {
			return nil, fmt.Errorf("error scanning product and category: %v", err)
		}

		if _, exists := productsMap[prodID]; !exists {
			productsMap[prodID] = &models.Product{
				ID:          prodID,
				Name:        prodName,
				Description: prodDescription,
				Price:       prodPrice,
				Categories:  []models.Category{},
			}
		}

		productsMap[prodID].Categories = append(productsMap[prodID].Categories, models.Category{
			ID:          catID,
			Name:        catName,
			Description: catDescription,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products and categories: %v", err)
	}

	products := make([]models.Product, 0, len(productsMap))
	for _, product := range productsMap {
		products = append(products, *product)
	}

	return products, nil
}

func UpdateProduct(productID int, updateReq models.ProductUpdateRequest) error {
    var setParts []string
    var args []interface{}
    var argIndex int = 1

    if updateReq.Name != nil {
        setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
        args = append(args, *updateReq.Name)
        argIndex++
    }
    if updateReq.Description != nil {
        setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
        args = append(args, *updateReq.Description)
        argIndex++
    }
    if updateReq.Price != nil {
        setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
        args = append(args, *updateReq.Price)
        argIndex++
    }

    if len(setParts) == 0 {
        return fmt.Errorf("no fields to update")
    }

    setClause := strings.Join(setParts, ", ")
    queryString := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", setClause, argIndex)
    args = append(args, productID)

    tx, err := db.Begin()
    if err != nil {
        return err
    }

    _, err = tx.Exec(queryString, args...)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("error updating product: %w", err)
    }

    _, err = tx.Exec(`DELETE FROM product_category WHERE product_id = $1`, productID)
    if err != nil {
        tx.Rollback()
        return ErrRollback
    }

    for _, categoryName := range updateReq.Categories {
        var categoryID int
        err := tx.QueryRow(`SELECT id FROM categories WHERE name = $1`, categoryName).Scan(&categoryID)
        if err != nil {
            tx.Rollback()
            return ErrCategoryDoesntExists
        }

        _, err = tx.Exec(`INSERT INTO product_category (product_id, category_id) VALUES ($1, $2)`, productID, categoryID)
        if err != nil {
            tx.Rollback()
            return ErrCreatingProduct
        }
    }

    if err := tx.Commit(); err != nil {
        tx.Rollback()
        return err
    }

    return nil
}

func DeleteProduct(productID int) error {
	tx, err := db.Begin()

	_, err = tx.Exec(`DELETE FROM product_category WHERE product_id = $1`, productID)
	if err != nil {
		tx.Rollback()
		return ErrRollback
	}

	_, err = tx.Exec("DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
