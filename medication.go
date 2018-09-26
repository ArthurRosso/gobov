package main

import "time"

type Medication struct {
	ID          int
	Description string
	Date        time.Time
	Animals     []Animal   `gorm:"many2many:medication_animal"`
	Medicines   []Medicine `gorm:"many2many:medication_medicine"`
	User        *User

	UserID int
}

func (m Medication) DateFmt() string {
	return m.Date.Format("02/01/2006")
}
