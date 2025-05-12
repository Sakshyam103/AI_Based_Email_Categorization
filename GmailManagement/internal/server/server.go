package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"GmailManagement/internal/auth"
	"GmailManagement/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
		fmt.Printf("WARNING: PORT environment variable not set or invalid, using default port: %d\n", port)
	} else {
		fmt.Printf("Using PORT from environment: %d\n", port)
	}

	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server will listen on http://localhost%s\n", addr)

	handler := NewServer.RegisterRoutes()

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	auth.NewAuth()

	return server
}
