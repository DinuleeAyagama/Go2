package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database connection
var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=nike port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate models
	db.AutoMigrate(&User{}, &Product{}, &Order{})
}

func main() {
	// Initialize database
	initDB()

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Product routes
	api.HandleFunc("/products", getProducts).Methods("GET")
	api.HandleFunc("/products/{id}", getProduct).Methods("GET")
	api.HandleFunc("/products", authMiddleware(createProduct)).Methods("POST")
	api.HandleFunc("/products/{id}", authMiddleware(updateProduct)).Methods("PUT")
	api.HandleFunc("/products/{id}", authMiddleware(deleteProduct)).Methods("DELETE")

	// User routes
	api.HandleFunc("/register", registerUser).Methods("POST")
	api.HandleFunc("/login", loginUser).Methods("POST")
	api.HandleFunc("/profile", authMiddleware(getUserProfile)).Methods("GET")

	// Order routes
	api.HandleFunc("/cart", authMiddleware(getCart)).Methods("GET")
	api.HandleFunc("/cart", authMiddleware(addToCart)).Methods("POST")
	api.HandleFunc("/checkout", authMiddleware(checkout)).Methods("POST")

	// Serve static files (for frontend)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server starting on :8000")
	log.Fatal(srv.ListenAndServe())
}