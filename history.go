package main

import "github.com/go-sql-driver/mysql"

type History struct {
	ID          int
	Description string
	Date        mysql.NullTime
	User        *User
	Animal      *Animal
	Medication  *Medication

	UserID       int
	AnimalID     int
	MedicationID int
}
