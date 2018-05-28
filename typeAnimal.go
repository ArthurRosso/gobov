package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
)

type TypeAnimal struct {
	ID          int
	Type        string
	Description string
	Animals     []Animal
}
