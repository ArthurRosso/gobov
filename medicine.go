package main

import "time"

type Medicine struct {
	ID          int
	Name        string
	Active      bool
	Expiration  time.Time
	Description string
	Type        *TypeMedicine
	Picture     []byte       `gorm:"type:byte"`
	Medications []Medication `gorm:"many2many:medication_medicine"`
	User        *User

	UserID int
	TypeID int
}

func NewMedicine() Medicine {
	return Medicine{Active: true}
}

func (m Medicine) ExpirationFmt() string {
	return m.Expiration.Format("02/01/2006")
}
