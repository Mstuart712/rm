package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mstuart712/rm/api/models"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error
	if Dbdriver == "postgres" {
		// DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		// server.DB, err = gorm.Open(Dbdriver, DBURL)
		user := getEnv("PG_USER", DbUser)
		password := getEnv("PG_PASSWORD", DbPassword)
		host := getEnv("PG_HOST", DbHost)
		port := getEnv("PG_PORT", DbPort)
		database := getEnv("PG_DB", DbName)

		dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			user,
			password,
			host,
			port,
			database,
		)
		server.DB, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Character{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
