package main

import (
	"log"
	"net/http"
	"os"
	"review-service/cache"
	DB "review-service/db"
	"review-service/handlers"
	"review-service/middlewares"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := cache.NewRedis(redisAddr)
	if err := cache.ConnectRedisWithRetry(rdb); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	publicKey, err := middlewares.LoadPublicKey("./keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := middlewares.AuthMiddleware(publicKey)

	h := handlers.New(db, rdb)
	fs := http.FileServer(http.Dir("./uploads/reviews"))
	mux := http.NewServeMux()
	mux.HandleFunc("/review/search", h.SearchReviewsHandler)
	mux.HandleFunc("/review/get", h.GetReviewHandler)
	mux.HandleFunc("/city/all", h.GetCitiesHandler)
	mux.HandleFunc("/city", h.GetCityByIDHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/review/create", authMiddleware(http.HandlerFunc(h.CreateReviewHandler)))
	mux.Handle("/review/like", authMiddleware(http.HandlerFunc(h.LikeReviewHandler)))
	mux.Handle("/review/liked", authMiddleware(http.HandlerFunc(h.GetLikesByUser)))
	mux.Handle("/review/delete", authMiddleware(http.HandlerFunc(h.DeleteReviewHandler)))
	mux.Handle("/review/status/update", authMiddleware(http.HandlerFunc(h.UpdateReviewStatusHandler)))
	handlerWithCors := middlewares.CorsMiddleware(mux)
	log.Println("Review service started on port 8080")
	log.Println(http.ListenAndServe(":8080", handlerWithCors))
}
