package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
)

type Breed struct {
	ID          int
	Breed       string
	Description string
	Animals     []Animal
}
