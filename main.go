package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *gorm.DB

func main() {
	// Retrieve environment variables
	dbHost := os.Getenv("DB_HOST")        // Should be 'mysql_container'
	dbUser := os.Getenv("DB_USER")        // Should be 'root'
	dbPassword := os.Getenv("DB_PASSWORD") // MySQL root password
	dbName := os.Getenv("DB_NAME")        // Database name

	// Build the DSN string for MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)

	// Open a connection to the MySQL database
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Ensure the 'users' table exists
	log.Println("Checking if the 'users' table exists...")
	if db.Migrator().HasTable(&User{}) {
		log.Println("'users' table already exists.")
	} else {
		log.Println("'users' table does not exist. Creating it now...")
		if err := db.AutoMigrate(&User{}); err != nil {
			log.Fatalf("Failed to create the 'users' table: %v", err)
		}
		log.Println("'users' table created successfully!")
	}

	// Initialize the router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start the server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// createUser handles POST /users
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// getUsers handles GET /users
func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// getUser handles GET /users/{id}
func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// updateUser handles PUT /users/{id}
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedData User
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.Name = updatedData.Name
	user.Email = updatedData.Email

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// deleteUser handles DELETE /users/{id}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&User{}, id).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

