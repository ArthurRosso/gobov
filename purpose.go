package main

type Purpose struct {
	ID          int
	Purpose     string
	Description string
	Animals     []Animal `gorm:"many2many:animal_purpose"`
	User        *User

	UserID int
}

func (p Purpose) String() string {
	return p.Purpose
}
