package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Cat struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Breed     string    `gorm:"size:255;not null;unique" json:"breed"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Cat) Prepare() {
	p.ID = 0
	p.Breed = html.EscapeString(strings.TrimSpace(p.Breed))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Cat) Validate() error {

	if p.Breed == "" {
		return errors.New("Required Breed")
	}
	return nil
}

func (p *Cat) FindCatByID(db *gorm.DB, pbreed string) (*Cat, error) {
	var err error
	err = db.Debug().Model(&Cat{}).Where("breed = ?", pbreed).Take(&p).Error
	if err != nil {
		return &Cat{}, err
	}
	return p, nil
}

func (p *Cat) FindAllCats(db *gorm.DB) (*[]Cat, error) {
	var err error
	cats := []Cat{}
	err = db.Debug().Model(&Cat{}).Limit(100).Find(&cats).Error
	if err != nil {
		return &[]Cat{}, err
	}
	if len(cats) > 0 {
		for i, _ := range cats {
			err := db.Debug().Model(&User{}).Where("id = ?", cats[i].ID).Take(&cats[i].Breed).Error
			if err != nil {
				return &[]Cat{}, err
			}
		}
	}
	return &cats, nil
}

func (p *Cat) SaveCat(db *gorm.DB) (*Cat, error) {
	var err error
	err = db.Debug().Model(&Cat{}).Create(&p).Error
	if err != nil {
		return &Cat{}, err
	}
	return p, nil
}
