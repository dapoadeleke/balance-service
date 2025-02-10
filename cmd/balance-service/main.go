package main

import (
	"github.com/dapoadeleke/balance-service/internal/db"
	"github.com/dapoadeleke/balance-service/internal/http"
	"github.com/dapoadeleke/balance-service/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbConn := db.NewPostgres(dbHost, dbUser, dbPassword, dbName)
	defer func(dbConn *db.Postgres) {
		err := dbConn.Close()
		if err != nil {
			log.WithError(err).Error("failed to close db connection")
		}
	}(dbConn)

	migrations.RunMigration(dbHost, dbUser, dbPassword, dbName)

	// Logger
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	handler := http.NewHandler(dbConn, logger)

	http.BuildRoutes(app, handler)

	logger.Fatal(app.Listen(":8080"))
}
