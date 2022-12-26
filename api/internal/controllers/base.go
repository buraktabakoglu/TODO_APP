package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	//"github.com/gorilla/mux"
	"github.com/gin-gonic/gin"

	_"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	//server.DB.Debug().AutoMigrate(&models.User{}, &models.Todo{}) //database migration yapılacak.
	server.DB.CreateTable(&models.User{})
	server.DB.CreateTable(&models.Todo{})
	

	server.Router = gin.Default()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
