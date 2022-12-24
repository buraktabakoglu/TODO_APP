package test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/buraktabakoglu/GOLANGAPPX/api/internal/controllers"
	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}


var userInstance = models.User{}

var todoInstance = models.Todo{}
	
var err error

func TestMain(m *testing.M) {

	

	
		
		err = godotenv.Load(os.ExpandEnv("../../.env"))
		if err != nil {
			log.Fatalf("error getting env %v\n", err)
		}
		Database()
		os.Exit(m.Run())
	}

		
	
	
	



func Database() {
	var err error
	

	TestDbDriver := os.Getenv("TestDbDriver")
	

	fmt.Println("st")

	if TestDbDriver == "postgres" {
		fmt.Println("stafasffsafsaasffsa")
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		fmt.Println(err)

		if err != nil {			
			fmt.Printf("cannot connect to %s database\n", TestDbDriver)
			log.Fatal("this is the error:", err)
		} else {
			fmt.Printf("we are connected to the %s database\n", TestDbDriver)
		}
	}
	
	
}

func refreshUserTable() error {
	
	
	
	fmt.Println(server)
	fmt.Println(server.DB)

	
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Nickname: "naber",
		Email:    "naber17@gmail.com",
		Password: "password17",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		{
			Nickname: "burak tabak",
			Email:    "burak12@gmail.com",
			Password: "password12",
		},
		{
			Nickname: "banshe tabak",
			Email:    "banshe13@gmail.com",
			Password: "password13",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndTodoTable() error {
	

	err := server.DB.DropTableIfExists(&models.User{}, &models.Todo{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Todo{}).Error
	if err != nil {
		return err
	}
	log.Printf("successfully refreshed tables")
	return nil
}

func seedOneUserAndOneTodo() (models.Todo, error) {
	

	err := refreshUserAndTodoTable()
	if err != nil {
		return models.Todo{}, err
	}
	user := models.User{
		Nickname: "hobbit",
		Email:    "hobbit@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Todo{}, err
	}
	todo := models.Todo{
		Status:"This is the title hobbit",
		Description:"This is the content hobbit",
		AuthorID: uint32(user.ID),
	}
	err = server.DB.Model(&models.Todo{}).Create(&todo).Error
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func seedUsersAndTodos() ([]models.User, []models.Todo, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Todo{}, err
	}
	var users = []models.User{
		{
			Nickname: "gandalf",
			Email:    "gandalf@gmail.com",
			Password: "password123",
		},
		{
			Nickname: "yildimamk",
			Email:    "yildim@gmail.com",
			Password: "password11",
		},
	}
	var todos = []models.Todo{
		{
			Status:"status 1",
			Description:"hello world 1",
		},
		{
			Status:"status 2",
			Description:"hello world 2",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		todos[i].AuthorID = uint32(users[i].ID)

		err = server.DB.Model(&models.Todo{}).Create(&todos[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, todos, nil
}
