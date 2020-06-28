package seed

import (
	"log"

	"../../api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "admin",
		Password: "@#$RF@!718",
	},
}

var cats = []models.Cat{
	models.Cat{
		Breed: "Ragdoll",
	},
	models.Cat{
		Breed: "Persa",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Cat{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Cat{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
	for i, _ := range cats {
		err = db.Debug().Model(&models.Cat{}).Create(&cats[i]).Error
		if err != nil {
			log.Fatalf("cannot seed cats table: %v", err)
		}
	}

}
