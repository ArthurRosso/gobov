package main

type User struct {
	ID          int
	Username    string
	Email       string
	Password    string
	Animals     []Animal
	Medicines   []Medicine
	Medications []Medication
}
