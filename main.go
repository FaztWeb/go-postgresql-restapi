package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Tasks []Task
}

type Task struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

var db *gorm.DB
var err error

func main() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	port := os.Getenv("PORT")

	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port)

	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Task{})

	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	http.ListenAndServe(":3000", r)

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []Person
	db.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user Person
	var tasks []Task

	db.First(&user, params["id"])
	db.Model(&user).Association("Tasks").Find(&tasks)

	user.Tasks = tasks
	json.NewEncoder(w).Encode(&user)
}
