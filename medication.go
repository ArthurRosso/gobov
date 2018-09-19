package main

import (
	"github.com/go-sql-driver/mysql"
)

type Medication struct {
	ID          int
	Description string
	Date        mysql.NullTime
	Animals     []Animal   `gorm:"many2many:medication_animal"`
	Medicines   []Medicine `gorm:"many2many:medication_medicine"`
}

func (m Medication) DateFmt() string {
	return m.Date.Time.Format("02/01/2006")
}
