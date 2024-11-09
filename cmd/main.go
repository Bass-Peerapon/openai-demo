package main

import (
	"log"
	"os"

	"github.com/Bass-Peerapon/openai-demo/internal/repository/postgres"
	"github.com/Bass-Peerapon/openai-demo/internal/repository/resty"
	"github.com/Bass-Peerapon/openai-demo/internal/rest"
	"github.com/Bass-Peerapon/openai-demo/openai"
	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var (
	APP_PORT          = os.Getenv("APP_PORT")
	CHAT_DATABASE     = os.Getenv("CHAT_DATABASE")
	CUSTOMER_DATABASE = os.Getenv("CUSTOMER_DATABASE")
	PRODUCT_DATABASE  = os.Getenv("PRODUCT_DATABASE")
	OPENAI_HOST       = os.Getenv("OPENAI_HOST")
	OPENAI_SECRET     = os.Getenv("OPENAI_SECRET")
	OPENAI_MODEL      = os.Getenv("OPENAI_MODEL")
)

func ConnectPostgres(conn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func runMigrations(db *sqlx.DB, path string) error {
	// Initialize migrate with a PostgreSQL database instance
	driver, err := migrate_postgres.WithInstance(db.DB, &migrate_postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,       // Path to the migrations folder
		"postgres", // Database name
		driver,
	)
	if err != nil {
		return err
	}

	// Apply the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	chatClinet := ConnectPostgres(CHAT_DATABASE)
	customerClient := ConnectPostgres(CUSTOMER_DATABASE)
	productClient := ConnectPostgres(PRODUCT_DATABASE)

	if err := runMigrations(chatClinet, "file://migrate/chat"); err != nil {
		log.Fatal(err)
	}
	if err := runMigrations(customerClient, "file://migrate/customer"); err != nil {
		log.Fatal(err)
	}
	if err := runMigrations(productClient, "file://migrate/product"); err != nil {
		log.Fatal(err)
	}

	// Prepare Repository
	chatRepo := postgres.NewChatRepository(chatClinet)
	customerRepo := postgres.NewCustomerRepository(customerClient)
	openaiRepo := resty.NewOpenaiService(OPENAI_HOST, OPENAI_SECRET, OPENAI_MODEL)
	productRepo := postgres.NewProductRepository(productClient, openaiRepo)

	if err := productRepo.MigrateData(); err != nil {
		log.Fatal(err)
	}

	// Prepare Service
	openaiService := openai.NewService(openaiRepo, customerRepo, productRepo, chatRepo)

	rest.NewOpenaiHandler(e, openaiService)

	log.Fatal(e.Start(APP_PORT))
}
