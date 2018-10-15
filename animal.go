package main

import (
	"fmt"
	"math"
	"time"

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
	return a.Birthday.Time.Format("02/01/2006")
}

func (a Animal) MotherFmt() string {
	if a.Mother == nil {
		return "Sem mãe"
	} else {
		return a.Mother.Name
	}
}

func (a Animal) FatherFmt() string {
	if a.Father == nil {
		return "Sem pai"
	} else {
		return a.Father.Name
	}
}

func (a Animal) BirthPFmt() string {
	return a.Birthday.Time.Format("2006-01-02")
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
	b := a.Birthday.Time
	x := RoundTime(b.Sub(time.Now()).Seconds() / 2600640)
	return x * (-1)
}

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}
