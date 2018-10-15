package main

type TypeAnimal struct {
	ID          int
	Type        string
	Description string
	Animals     []Animal
	User        *User

	UserID int
}
