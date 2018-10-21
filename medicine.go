package main

import "time"

type Medicine struct {
	ID          int
	Name        string
	Expiration  time.Time
	Description string
	Type        *TypeMedicine
	Picture     []byte       `gorm:"type:bytea"`
	Medications []Medication `gorm:"many2many:medication_medicine"`
	User        *User

	UserID int
	TypeID int
}

func NewMedicine() Medicine {
	return Medicine{}
}

func (m Medicine) ExpirationFmt() string {
	return m.Expiration.Format("02/01/2006")
}

func (m Medicine) ExpFmt() string {
	return m.Expiration.Format("2006-01-02")
}
