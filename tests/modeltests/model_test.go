package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"../../api/controllers"
	"../../api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var catInstance = models.Cat{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Nickname: "Pet",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		models.User{
			Nickname: "Steven victor",
			Password: "password",
		},
		models.User{
			Nickname: "Kenny Morris",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndCatTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Cat{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Cat{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedUsersAndCats() ([]models.User, []models.Cat, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Cat{}, err
	}
	var users = []models.User{
		models.User{
			Nickname: "Steven victor",
			Password: "password",
		},
		models.User{
			Nickname: "Magu Frank",
			Password: "password",
		},
	}
	var cats = []models.Cat{
		models.Cat{
			Breed: "Breed 1",
		},
		models.Cat{
			Breed: "Breed 2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		cats[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Cat{}).Create(&cats[i]).Error
		if err != nil {
			log.Fatalf("cannot seed cats table: %v", err)
		}
	}
	return users, cats, nil
}
