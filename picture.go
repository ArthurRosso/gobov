package main

type Picture struct {
	ID      int
	Main    bool
	Picture []byte `gorm:"type:byte"`
	Animal  *Animal

	AnimalID int
}
