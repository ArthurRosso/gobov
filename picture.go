package main

type Picture struct {
	ID      int
	Main    bool
	Picture []byte `gorm:"type:mediumblob"`
	Animal  *Animal

	AnimalID int
}
