package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres

	"github.com/gotutorial/api/middlewares"
	"github.com/gotutorial/api/models"
	"github.com/gotutorial/api/responses"
)

type App struct {
    Router *mux.Router
    DB     *gorm.DB
}

// Initialize connect to the database and wire up routes
func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
    var err error
    DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

    a.DB, err = gorm.Open("postgres", DBURI)
    if err != nil {
        fmt.Printf("\n Cannot connect to database %s", DbName)
        log.Fatal("This is the error:", err)
    } else {
        fmt.Printf("We are connected to the database %s", DbName)
    }

    a.DB.Debug().AutoMigrate(&models.User{}, &models.Task{}, &models.PendingTask{}) //database migration

    a.Router = mux.NewRouter().StrictSlash(true)
    a.initializeRoutes()
}

func (a *App) initializeRoutes() {
    a.Router.Use(middlewares.SetContentTypeMiddleware) // setting content-type to json

    a.Router.HandleFunc("/", home).Methods("GET")
    a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
    a.Router.HandleFunc("/login", a.Login).Methods("POST")

    s := a.Router.PathPrefix("/api").Subrouter() // routes that require authentication
    s.Use(middlewares.AuthJwtVerify)

    s.HandleFunc("/tasks", a.CreateTask).Methods("POST")
    s.HandleFunc("/tasks", a.GetTasks).Methods("GET")
    s.HandleFunc("/tasks/{id:[0-9]+}", a.UpdateTask).Methods("PUT")
    s.HandleFunc("/tasks/{id:[0-9]+}", a.DeleteTask).Methods("DELETE")
    s.HandleFunc("/assignTask/{id:[0-9]+}", a.AssignTask).Methods("POST")
}

func (a *App) RunServer() {
    log.Printf("\nServer starting on port 5000")
    log.Fatal(http.ListenAndServe(":5000", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
    responses.JSON(w, http.StatusOK, "Welcome To Go Tutorial")
}