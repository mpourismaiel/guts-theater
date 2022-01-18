package config

import "os"

type Config struct {
	Address    string
	Port       string
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	TestMode   bool
}

func Setup() *Config {
	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "admin"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "guts"
	}

	return &Config{
		Address:    address,
		Port:       port,
		DbHost:     dbHost,
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		TestMode:   false,
	}
}
