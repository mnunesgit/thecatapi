package modeltests

import (
	"log"
	"testing"

	"../../api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllCats(t *testing.T) {

	err := refreshUserAndCatTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndCats()
	if err != nil {
		log.Fatalf("Error seeding user and post  table %v\n", err)
	}
	posts, err := catInstance.FindAllCats(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestSaveCat(t *testing.T) {

	err := refreshUserAndCatTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newCat := models.Cat{
		ID:    1,
		Breed: "This is the breed",
	}
	savedCat, err := newCat.SaveCat(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newCat.ID, savedCat.ID)
	assert.Equal(t, newCat.Breed, savedCat.Breed)

}

func TestGetCatByID(t *testing.T) {

	err := refreshUserAndCatTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	cat, err := seedOneUserAndOneCat()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundCat, err := catInstance.FindCatByID(server.DB, cat.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundCat.ID, cat.ID)
	assert.Equal(t, foundCat.Breed, cat.Breed)
}
