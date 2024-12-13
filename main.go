package main

import (
	"database/sql"
	"fmt"
	env "password_generator/package/environment"
	"password_generator/repository"

	docker "password_generator/internal/container_manager"
	flags "password_generator/internal/flag_manager"
	types "password_generator/internal/types"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Loading .env variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	dbOptions := types.DatabaseOptions{
		ContainerName: env.GetEnv("CONTAINER_NAME"),
		UserName:      env.GetEnv("CONTAINER_USER"),
		DbPassword:    env.GetEnv("DB_PASSWORD"),
		DbPort:        env.GetEnv("DB_PORT"),
		DbName:        env.GetEnv("DB_NAME"),
		DbHost:        env.GetEnv("DB_HOST"),
		SslMode:       env.GetEnv("SSL_MODE"),
		HostPort:      env.GetEnv("HOST_PORT"),
	}

	// Managing a docker container before working with DB
	docker.CreateAndStartContainer(dbOptions)

	// Connecting to database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbOptions.UserName, dbOptions.DbPassword, dbOptions.DbName, dbOptions.DbHost, dbOptions.DbPort, dbOptions.SslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("Failing connecting to database: %v", err))
	}
	defer db.Close()

	// Postgres connection check
	repository.CheckConnection(db)
	repository.CreateTable(db)

	// Console flags management
	flags.ManageFlags(db)
}
