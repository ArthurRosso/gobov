package main

import (
	"github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
)

type Weight struct {
	ID          int
	Weight      float32
	Description string
	Date        mysql.NullTime
	Animal      *Animal

	AnimalID 	int
}

func (w Weight) DateFmt() string {
	return w.Date.Time.Format("02/01/2006")
}