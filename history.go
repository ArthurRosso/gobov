package main

import "time"

type History struct {
	ID          int
	Description string
	Date        time.Time
	User        *User
	Animals     []*Animal `gorm:"many2many:animal_history"`
	Medication  *Medication

	UserID       int
	AnimalID     int
	MedicationID int
}
