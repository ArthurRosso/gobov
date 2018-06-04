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
	db.AutoMigrate(&Medicine{})
	db.AutoMigrate(&TypeMedicine{})
	db.AutoMigrate(&Picture{})

	Data()
	DataM()

	r := mux.NewRouter()

	r.HandleFunc("/pic/{idAnimal}", getPic)

	r.HandleFunc("/animal", getAnimal)
	r.HandleFunc("/newAnimal", postAnimal)
	r.HandleFunc("/medicine", getMedicine)
	r.HandleFunc("/newMedicine", postMedicine)
	r.HandleFunc("/profile/{idAnimal}", getProfile)
	http.ListenAndServe(":8000", r)
}

func getPic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	picture := Picture{}
	db.First(&picture, idAnimal)
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

	r.ParseMultipartForm(0)
	m := r.MultipartForm

	files := m.File["Pictures"]

	db.Save(&animal)

	if len(files) > 0 {
		pic := Picture{Main: true, AnimalID: animal.ID}
		arquivo, _ := files[0].Open()
		pic.Picture, _ = ioutil.ReadAll(arquivo)
		defer arquivo.Close()
		db.Save(&pic)
		for _, file := range files[1:] {
			picture := Picture{AnimalID: animal.ID}
			arquivo, _ := file.Open()
			picture.Picture, _ = ioutil.ReadAll(arquivo)
			defer arquivo.Close()
			db.Save(&picture)
		}
	}

	http.Redirect(w, r, "/animal", http.StatusMovedPermanently)
}

func getMedicine(w http.ResponseWriter, r *http.Request) {

	medicines := []Medicine{}
	db.Preload("Type").Find(&medicines, Medicine{})

	types := []TypeMedicine{}
	db.Find(&types, &TypeMedicine{})

	context := map[string]interface{}{
		"types":     types,
		"medicines": medicines,
	}

	str, _ := mustache.RenderFile("templates/medicine.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postMedicine(w http.ResponseWriter, r *http.Request) {
	medicine := NewMedicine()
	medicine.Name = r.PostFormValue("Name")

	expiration, _ := time.Parse("2006-01-02", r.PostFormValue("Expiration"))
	medicine.Expiration = mysql.NullTime{Time: expiration, Valid: true}

	medicine.Description = r.PostFormValue("Description")

	typeM := TypeMedicine{}
	idType, _ := strconv.Atoi(r.PostFormValue("Type"))
	db.Find(&typeM, idType)
	medicine.Type = &typeM
	db.First(&medicine.Type, idType)

	db.Save(&medicine)

	http.Redirect(w, r, "/medicine", http.StatusMovedPermanently)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{}
	db.Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").First(&animal, idAnimal)

	context := map[string]interface{}{
		"animal": animal,
	}

	str, _ := mustache.RenderFile("templates/profile.html", context)
	bit := []byte(str)
	w.Write(bit)
}
