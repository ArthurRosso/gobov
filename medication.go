package main

import (
	"fmt"
	"time"
)

type Medication struct {
	ID          int
	Description string
	Date        time.Time
	Animals     []*Animal  `gorm:"many2many:medication_animal"`
	Medicines   []Medicine `gorm:"many2many:medication_medicine"`
	User        *User

	UserID int
}

func (m Medication) DateFmt() string {
	return m.Date.Format("02/01/2006")
}

func (m Medication) AnimalsFmt() string {
	var res string
	a := m.Animals[0]
	res += fmt.Sprint(a.Name)
	for _, element := range m.Animals[1:] {
		res += ", " + fmt.Sprint(element.Name)
	}
	return res
}

func (m Medication) MedicinesFmt() string {
	var res string
	me := m.Medicines[0]
	res += fmt.Sprint(me.Name)
	for _, element := range m.Medicines[1:] {
		res += ", " + fmt.Sprint(element.Name)
	}
	return res
}
