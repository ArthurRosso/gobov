package main

type Picture struct {
	ID      int
	Main    bool
	Picture []byte `gorm:"type:bytea"`
	Animal  *Animal

	AnimalID int
}
