package main

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

type Weight struct {
	ID          int
	Weight      float32
	Description string
	Date        mysql.NullTime
	Animal      *Animal

	AnimalID int
}

func (w Weight) DateFmt() string {
	return w.Date.Time.Format("02/01/2006")
}

func (w Weight) Months() time.Month {
	return w.Date.Time.Month()
}

func (w Weight) Years() int {
	return w.Date.Time.Year()
}
