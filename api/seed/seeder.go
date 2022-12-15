package seed

import (
	"log"

	"github.com/buraktabakoglu/GOLANGAPPX/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Nickname: "Burak tabakoğlu",
		Email:    "burakt1@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Visual Code",
		Email:    "vs@gmail.com",
		Password: "password",
	},
}

var todos = []models.Todo{
	{

		Status:      "online",
		Description: "Vs golangapp 1",
	},
	{

		Status:      "ofline",
		Description: "Vs golangapp 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Todo{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Todo{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	/*err = db.Debug().Model(&models.Todo{}).AddForeignKey("status", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}*/

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		todos[i].ID = users[i].ID

		err = db.Debug().Model(&models.Todo{}).Create(&todos[i]).Error
		if err != nil {
			log.Fatalf("cannot seed todos table: %v", err)
		}
	}
}
