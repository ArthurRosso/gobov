package main

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type Animal struct {
	ID          int
	Name        string
	Active      bool
	Birthday    mysql.NullTime
	Weights     []Weight
	Type        *TypeAnimal
	Breed       *Breed
	Purposes    []Purpose `gorm:"many2many:animal_purpose"`
	Pictures    []Picture
	Medications []Medication `gorm:"many2many:medication_animal"`
	Mother      *Animal
	Father      *Animal
	Sons        []Animal

	MotherID int
	FatherID int
	TypeID   int
	BreedID  int
}

func NewAnimal() Animal {
	return Animal{Active: true}
}

func (a Animal) BirthFmt() string {
	return a.Birthday.Time.Format("02/01/2006")
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
	db.Where("main=?", 1).First(&pic)
	return pic
}
