package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Animal struct {
	ID          int
	Name        string
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

func (a Animal) Children() int {
	var count int
	animals := []Animal{}
	db.Find(&animals, Animal{})

	db.Where("mother_id = ?", a.ID).Or("father_id = ?", a.ID).Find(&animals).Count(&count)
	//// SELECT * from USERS WHERE name = 'jinzhu' OR name = 'jinzhu 2'; (users)
	//// SELECT count(*) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'; (count)
	return count
}

func NewAnimal() Animal {
	return Animal{}
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

func (a Animal) MainPicFmt() int {
	pic := Picture{}
	db.Where("main=? and animal_id=?", 1, a.ID).First(&pic)
	return pic.ID
}

func (a Animal) MotherPicFmt() string {
	if a.Mother == nil {
		return ""
	} else {
		return "/pic/" + strconv.Itoa(a.Mother.MainPicFmt())
	}
}

func (a Animal) FatherFmt() string {
	if a.Father == nil {
		return "Sem pai"
	} else {
		return a.Father.Name
	}
}

func (a Animal) FatherPicFmt() string {
	if a.Father == nil {
		return ""
	} else {
		return "/pic/" + strconv.Itoa(a.Father.MainPicFmt())
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

func (a Animal) WeightPFmt() float32 {
	if len(a.Weights) > 0 {
		return a.Weights[len(a.Weights)-1].Weight
	} else {
		return 0
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

func (a Animal) Age() string {
	// Your birthday: let's say it's November 6st, 2000, 0:00 AM
	birthday := a.Birthday.Time
	year, month, day, hour, min, sec := diff(birthday, time.Now())

	fmt.Printf("You are %d years, %d months, %d days, %d hours, %d mins and %d seconds old.", year, month, day, hour, min, sec)

	return strconv.Itoa(year) + " anos e " + strconv.Itoa(month) + " meses"
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
