package api

import (
	"fmt"
	"log"
	"os"

	"github.com/buraktabakoglu/GOLANGAPPX/api/internal/controllers"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var err error

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("env.file found")
	}
}

func Run() {

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("error getting env, not comming through %v", err)
	} else {
		fmt.Println("we are getting the env value")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))


	// Start server
	server.Run(":8080")

}
