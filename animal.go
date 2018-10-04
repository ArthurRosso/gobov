package main

import (
	"fmt"
	"time"
)

type Animal struct {
	ID          int
	Name        string
	Active      bool
	Birthday    time.Time
	Weights     []Weight
	Type        *TypeAnimal
	Breed       *Breed
	Purposes    []Purpose `gorm:"many2many:animal_purpose"`
	Pictures    []Picture
	Medications []Medication `gorm:"many2many:medication_animal"`
	Mother      *Animal
	Father      *Animal
	User        *User
	Histories   []History `gorm:"many2many:animal_history"`

	UserID   int
	MotherID int
	FatherID int
	TypeID   int
	BreedID  int
}

func NewAnimal() Animal {
	return Animal{Active: true}
}

func (a Animal) BirthFmt() string {
	return a.Birthday.Format("02/01/2006")
}

func (a Animal) BirthPFmt() string {
	return a.Birthday.Format("2006-01-02")
}

func (a Animal) WeightFmt() string {
	if len(a.Weights) > 0 {
		return fmt.Sprint(a.Weights[len(a.Weights)-1].Weight)
	} else {
		return ":("
	}
}

func (a Animal) PurposesFmt() string {
	var res string
	if len(a.Purposes) > 0 {
		res += fmt.Sprint(a.Purposes[0])
		for _, element := range a.Purposes[1:] {
			res += ", " + fmt.Sprint(element.Purpose)
		}
		return res
	} else {
		return ":("
	}
}

func (a Animal) MainPic() Picture {
	pic := Picture{}
	db.Where("main=? and animal_id=?", 1, a.ID).First(&pic)
	return pic
}

func (a Animal) Age() int {
	now := time.Now()
	years := now.Year() - a.Birthday.Year()
	if now.YearDay() < a.Birthday.YearDay() {
		years--
	}
	return years
}
