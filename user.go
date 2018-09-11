package main

import (
)

type User struct {
	ID			int
	Username	string
	Email		string
	Name 		string
	Password	string
	Animals    	[]Animal
	Medicines   []Medicine
	Medications []Medication
}