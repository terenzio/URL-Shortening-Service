package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/terenzio/URL-Shortening-Service/application"
	"github.com/terenzio/URL-Shortening-Service/docs"
	urlHandler "github.com/terenzio/URL-Shortening-Service/infrastructure/http"
	redisRepo "github.com/terenzio/URL-Shortening-Service/infrastructure/redis"
)

// @title          URL Shortening Service
// @version        1.0
// @description    This is a URL shortening service that allows users to shorten long URLs especially built for TSMC.
// @termsOfService http://swagger.io/terms/
// @contact.name   Terence Liu
// @contact.url    https://github.com/terenzio/URL-Shortening-Service
// @contact.email  terenzio@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:9000
// @BasePath  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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

	// Initialize the Gin router
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Register routes with their handlers
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := router.Group("/api/v1")
	{
		urlPage := v1.Group("/url")
		{
			urlPage.GET("/display", handler.HandleHomePage)
			urlPage.POST("/add", handler.HandleAddLink)
		}
		urlRedirect := v1.Group("/redirect")
		{
			urlRedirect.GET("/:shortcode", handler.HandleRedirectToOriginalLink)
		}
	}

	log.Println("\nThe URL Shortening Service is now running!")
	if err := router.Run(":9000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
