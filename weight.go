package main

import (
	"time"
)

type Weight struct {
	ID          int
	Weight      float32
	Description string
	Date        time.Time
	Animal      *Animal

	AnimalID int
}

func (w Weight) DateFmt() string {
	return w.Date.Format("02/01/2006")
}

func (w Weight) Months() time.Month {
	return w.Date.Month()
}

func (w Weight) Years() int {
	return w.Date.Year()
}
