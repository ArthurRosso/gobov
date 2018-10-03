package main

import "github.com/go-sql-driver/mysql"

type Medicine struct {
	ID          int
	Name        string
	Active      bool
	Expiration  mysql.NullTime
	Description string
	Type        *TypeMedicine
	Picture     []byte       `gorm:"type:bytea"`
	Medications []Medication `gorm:"many2many:medication_medicine"`
	User        *User

	UserID int
	TypeID int
}

func NewMedicine() Medicine {
	return Medicine{Active: true}
}

func (m Medicine) ExpirationFmt() string {
	return m.Expiration.Time.Format("02/01/2006")
}

func (m Medicine) ExpFmt() string {
	return m.Expiration.Time.Format("2006-01-02")
}
