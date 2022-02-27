package models

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type Character struct {
	ID     uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name   string `gorm:"size:100;not null;unique" json:"name"`
	Level  uint64 `gorm:"not null;" json:"level"`
	RaceId uint64 `json:"raceId"`
}

func (p *Character) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Level = 0
	p.RaceId = 0
}

func (p *Character) Validate() error {

	if p.Name == "" {
		return errors.New("Required Name")
	}
	if p.Level == 0 {
		return errors.New("Required Level")
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
		err = db.Debug().Model(&User{}).Where("id = ?", c.Name).Take(&c.Name).Error
		if err != nil {
			return &Character{}, err
		}
	}
	return c, nil
}

func (c *Character) FindCharacterByID(db *gorm.DB, pid uint64) (*Character, error) {
	var err error
	err = db.Debug().Model(&Character{}).Where("id = ?", pid).Take(&c).Error
	if err != nil {
		return &Character{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.Name).Take(&c.Name).Error
		if err != nil {
			return &Character{}, err
		}
	}
	return c, nil
}
