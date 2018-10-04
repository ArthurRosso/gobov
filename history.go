package main

import "github.com/go-sql-driver/mysql"

type History struct {
	ID          int
	Description string
	Date        mysql.NullTime
	User        *User
	Animals     []*Animal `gorm:"many2many:animal_history"`
	Medication  *Medication

	UserID       int
	AnimalID     int
	MedicationID int
}
