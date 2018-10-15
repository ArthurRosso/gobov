package main

type Breed struct {
	ID          int
	Breed       string
	Description string
	Animals     []Animal
	User        *User

	UserID int
}
