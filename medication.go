package main

import (
	"github.com/go-sql-driver/mysql"
)

type Medicine struct {
	ID          int
	Name        string
	Date        mysql.NullTime
	Description string
	Animals     []Animal `gorm:"many2many:medication_animal"`
	Medicines   []Animal `gorm:"many2many:medication_medicine"`

	TypeID int
}

func (m Medicine) DateFmt() string {
	return m.Expiration.Time.Format("02/01/2006")
}
