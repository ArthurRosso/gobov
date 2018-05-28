package main

import (
	"github.com/cbroglie/mustache"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	db *gorm.DB
)

func main() {
	db, _ = gorm.Open("mysql", "goboi:goboi@/goboi")

	db.AutoMigrate(&Animal{})
	db.AutoMigrate(&Weight{})
	db.AutoMigrate(&TypeAnimal{})
	db.AutoMigrate(&Breed{})
	db.AutoMigrate(&Purpose{})
	db.AutoMigrate(&Picture{})

	r := mux.NewRouter()

	r.HandleFunc("/pic/{idAnimal}", getPic)

	r.HandleFunc("/animal", getAnimal)
	r.HandleFunc("/newAnimal", postAnimal)
	http.ListenAndServe(":8000", r)
}

func getPic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{}
	picture := Picture{}
	db.Find(&picture, idAnimal)
	if len(picture.Picture) > 0 {
		w.Write(picture.Picture)
	}
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	animals := []Animal{}
	db.Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Find(&animals, Animal{})

	types := []TypeAnimal{}
	db.Find(&types, &TypeAnimal{})

	breeds := []Breed{}
	db.Find(&breeds, &Breed{})

	purposes := []Purpose{}
	db.Find(&purposes, &Purpose{})

	context := map[string]interface{}{
		"types":    types,
		"breeds":   breeds,
		"purposes": purposes,
		"animals":  animals,
	}

	str, _ := mustache.RenderFile("templates/animal.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postAnimal(w http.ResponseWriter, r *http.Request) {
	animal := NewAnimal()
	weight := Weight{Description: "Primeira pesagem"}
	animal.Name = r.PostFormValue("Name")

	birth, _ := time.Parse("2006-01-02", r.PostFormValue("Birthday"))
	animal.Birthday = mysql.NullTime{Time: birth, Valid: true}

	peso, _ := strconv.ParseFloat(r.PostFormValue("Weight"), 32)
	weight.Weight = float32(peso)
	animal.Weights = append(animal.Weights, weight)

	typeA := TypeAnimal{}
	idType, _ := strconv.Atoi(r.PostFormValue("Type"))
	db.Find(&typeA, idType)
	animal.Type = &typeA
	db.First(&animal.Type, idType)

	breed := Breed{}
	idBreed, _ := strconv.Atoi(r.PostFormValue("Breed"))
	db.Find(&breed, idBreed)
	animal.Breed = &breed
	db.First(&animal.Breed, idBreed)

	r.ParseForm()
	for _, idPurposes := range r.Form["Purpose"] {
		// element is the element from someSlice for where we are
		purpose := Purpose{}
		id, _ := strconv.Atoi(idPurposes)
		db.Find(&purpose, id)
		animal.Purposes = append(animal.Purposes, purpose)
	}
	m := r.MultipartForm

	files := m.File("Picture")
	defer file.Close()

	picture := Picture{}
	picture.Picture, _ = ioutil.ReadAll(file)
	animal.Pictures = append(animal.Pictures, picture)

	db.Save(&animal)

	http.Redirect(w, r, "/animal", http.StatusMovedPermanently)
}