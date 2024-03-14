package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"

	"github.com/terenzio/URL-Shortening-Service/application"
	urlHandler "github.com/terenzio/URL-Shortening-Service/infrastructure/http"
	redisRepo "github.com/terenzio/URL-Shortening-Service/infrastructure/redis"
)

func main() {

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Create a new URL repository
	repo := redisRepo.NewURLRepository(rdb)

	// Create a new URL service
	service := application.NewURLService(repo)

	// Create a new URL handler
	handler := urlHandler.NewHandler(service)

	// Register the URL handler and start the server
	http.HandleFunc("/", handler.HandleHomePage)
	http.HandleFunc("/addLink", handler.HandleAddLink)
	http.HandleFunc("/short/", handler.HandleRedirectToOriginalLink)

	log.Println("Server is running on port :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
