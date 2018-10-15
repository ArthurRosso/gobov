package main

type TypeMedicine struct {
	ID          int
	Type        string
	Description string
	Medicines   []Medicine
	User        *User

	UserID int
}
