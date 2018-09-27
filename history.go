package main

import "time"

type History struct {
	ID          int
	Description string
	Date        time.Time
	User        *User
	Animal      *Animal
	Medication  *Medication

	UserID       int
	AnimalID     int
	MedicationID int
}
