package main

import (
	"github.com/go-sql-driver/mysql"
)

type Medicine struct {
	ID          int
	Name        string
	Active      bool
	Expiration  mysql.NullTime
	Description string
	Type        *TypeMedicine
	Picture     []byte `gorm:"type:mediumblob"`
	Medications []Medication `gorm:"many2many:medication_medicine"`

	TypeID int
}

func NewMedicine() Medicine {
	return Medicine{Active: true}
}

func (m Medicine) ExpirationFmt() string {
	return m.Expiration.Time.Format("02/01/2006")
}
