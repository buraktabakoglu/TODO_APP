package api

import (
	"fmt"
	"log"
	"os"

	"github.com/buraktabakoglu/GOLANGAPPX/api/controllers"
	"github.com/buraktabakoglu/GOLANGAPPX/api/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	
	
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("error getting env, not comming through %v", err)
	} else {
		fmt.Println("we are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")

}
