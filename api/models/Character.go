package models

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type Character struct {
	ID      uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name    string `gorm:"size:100;not null;unique" json:"name"`
	Level   uint64 `gorm:"not null;" json:"level"`
	RaceId  uint64 `json:"raceId"`
	OwnerID uint32 `gorm:"not null" json:"owner_id"`
}

func (c *Character) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
}

func (c *Character) Validate() error {
	if c.Name == "" {
		return errors.New("Required Name")
	}
	if c.Level == 0 {
		return errors.New("Required Level")
	}
	if c.OwnerID < 1 {
		return errors.New("Required Owner")
	}
	return nil
}

func (c *Character) SaveCharacter(db *gorm.DB) (*Character, error) {
	var err error
	err = db.Debug().Model(&Character{}).Create(&c).Error
	if err != nil {
		return &Character{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.OwnerID).Take(&c.OwnerID).Error
		if err != nil {
			return &Character{}, err
		}
	}
	return c, nil
}

func (p *Character) FindAllCharacters(db *gorm.DB) (*[]Character, error) {
	var err error
	characters := []Character{}
	err = db.Debug().Model(&Character{}).Limit(100).Find(&characters).Error
	if err != nil {
		return &[]Character{}, err
	}
	// if len(characters) > 0 {
	// 	for i, _ := range characters {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", characters[i].OwnerID).Take(&characters[i].OwnerID).Error
	// 		if err != nil {
	// 			return &[]Character{}, err
	// 		}
	// 	}
	// }
	return &characters, nil
}

func (c *Character) FindCharacterByID(db *gorm.DB, pid uint64) (*Character, error) {
	var err error
	err = db.Debug().Model(&Character{}).Where("id = ?", pid).Take(&c).Error
	if err != nil {
		return &Character{}, err
	}
	// if c.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", c.OwnerID).Take(&c.Name).Error
	// 	if err != nil {
	// 		return &Character{}, err
	// 	}
	// }
	return c, nil
}

func (c *Character) UpdateACharacter(db *gorm.DB) (*Character, error) {

	var err error

	err = db.Debug().Model(&Character{}).Where("id = ?", c.ID).Updates(Character{Name: c.Name, RaceId: c.RaceId, Level: c.Level}).Error
	if err != nil {
		return &Character{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.OwnerID).Take(&c.OwnerID).Error
		if err != nil {
			return &Character{}, err
		}
	}
	return c, nil
}

func (p *Character) DeleteACharacter(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Character{}).Where("id = ? and owner_id = ?", pid, uid).Take(&Character{}).Delete(&Character{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
