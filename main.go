package main

import (
	"log"
	"net/http"
	"os"
	"review-service/auth"
	DB "review-service/db"
	"review-service/handlers"
)

func main() {
	cnf := DB.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := DB.ConnectWithRetry(cnf)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err := auth.LoadPublicKey("./keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := auth.AuthMiddleware(publicKey)

	h := handlers.New(db)
	mux := http.NewServeMux()
	log.Println("Auth service started on port 8081")
	log.Println(http.ListenAndServe(":8081", mux))
}
