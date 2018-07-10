package main

import (
	"fmt"
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
	db.AutoMigrate(&Medication{})

	r := mux.NewRouter()

	r.HandleFunc("/pic/{idAnimal}", getPic)
	r.HandleFunc("/picMedicine/{idMedicine}", getMedicinePic)

	r.HandleFunc("/animal", getAnimal)
	r.HandleFunc("/newAnimal", postAnimal)
	r.HandleFunc("/delAnimal/{ID}", delAnimal)
	/*
	r.HandleFunc("/editAnimal/{ID}", editAnimal)
	*/
	
	/*
	r.HandleFunc("/weight", getWeight)
	r.HandleFunc("/newWeight", postWeight)
	r.HandleFunc("/delWeight/{ID}", delWeight)
	*/

	r.HandleFunc("/medicine", getMedicine)
	r.HandleFunc("/newMedicine", postMedicine)
	r.HandleFunc("/delMedicine/{ID}", delMedicine)
	
	r.HandleFunc("/profile/{idAnimal}", getProfile)

	r.HandleFunc("/medication", getMedication)
	r.HandleFunc("/newMedication", postMedication)
	
	fmt.Println("Server listen and serve on port 8000")

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

func getProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{ID: idAnimal}
	db.Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Preload("Father").Preload("Mother").Preload("Pictures").First(&animal, idAnimal)

	context := map[string]interface{}{
		"animal": animal,
	}

	str, _ := mustache.RenderFile("templates/profile.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func getMedicinePic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	medicine := Medicine{}
	idMedicine, _ := strconv.Atoi(vars["idMedicine"])
	db.First(&medicine, idMedicine)
	if len(medicine.Picture) > 0 {
		w.Write(medicine.Picture)
	}
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	animals := []Animal{}
	db.Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Find(&animals, Animal{})

	fathers := []Animal{}
	mothers := []Animal{}
	for _, i := range animals {
		// element is the element from someSlice for where we are
		if i.Type.Type == "Touro" {
			fathers = append(fathers, i)
		} else if i.Type.Type == "Vaca" {
			mothers = append(mothers, i)
		}
	}

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
		"mothers":  mothers,
		"fathers":  fathers,
	}

	str, _ := mustache.RenderFile("templates/animal.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postAnimal(w http.ResponseWriter, r *http.Request) {
	animal := NewAnimal()
	animal.Name = r.PostFormValue("Name")

	motherID, _ := strconv.Atoi(r.PostFormValue("Mother"))
	if motherID != 0 {
		mother := Animal{}
		db.Find(&mother, Animal{ID: motherID})
		animal.Mother = &mother
	}

	fatherID, _ := strconv.Atoi(r.PostFormValue("Father"))
	if fatherID != 0 {
		father := Animal{}
		db.First(&father, Animal{ID: fatherID})
		animal.Father = &father
	}

	weight := Weight{Description: "Primeira pesagem"}

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

func delAnimal(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	id := m["ID"]

	// TODO: Delete Weights
	// TODO: Ver a linha 291

	//db.Model(&Medication{}).Where("animal_id = ?", id).Association("medication_animal").Delete(&Medication{})
	//db.Model(&Purpose{}).Where("animal_id = ?", id).Association("animal_purpose").Delete()// TODO: O que colocar e como
	db.Where("animal_id = ?", id).Delete(&Picture{})
	db.Where("ID = ?", id).Delete(&Animal{})

	http.Redirect(w, r, "/animal", http.StatusMovedPermanently)
}

// TODO: Ñão testado
/*
func editAnimal(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	id := m["ID"]

	animal.Name = r.PostFormValue("Name")

	motherID, _ := strconv.Atoi(r.PostFormValue("Mother"))
	if motherID != 0 {
		mother := Animal{}
		db.Find(&mother, Animal{ID: motherID})
		animal.Mother = &mother
	}

	fatherID, _ := strconv.Atoi(r.PostFormValue("Father"))
	if fatherID != 0 {
		father := Animal{}
		db.First(&father, Animal{ID: fatherID})
		animal.Father = &father
	}

	weight := Weight{Description: "Primeira pesagem"}

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

	http.Redirect(w, r, "/profile/{ID}", http.StatusMovedPermanently)
}
*/

func delMedicine(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	id := m["ID"]
	// TODO: A próxima linha não funciona, então possívelmente a linha 202 também não
	db.Model(&Medication{}).Where("medicine_id = ?", id).Association("medication_medicine").Delete(&Medication{})
	db.Where("ID = ?", id).Delete(&Medicine{})

	http.Redirect(w, r, "/medicine", http.StatusMovedPermanently)
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

	r.ParseMultipartForm(0)
	f := r.MultipartForm
	if f == nil {
		fmt.Println("Erro no formulário")
	}
	file := f.File["Picture"]
	arquivo, _ := file[0].Open()
	medicine.Picture, _ = ioutil.ReadAll(arquivo)
	defer arquivo.Close()

	db.Save(&medicine)

	http.Redirect(w, r, "/medicine", http.StatusMovedPermanently)
}

// TODO: Não testado
/*
func getWeight(w http.ResponseWriter, r *http.Request) {
	weights := []Weight{}
	db.Find(&weights, Weight{})

	context := map[string]interface{}{
		"weights":     weights,
	}

	str, _ := mustache.RenderFile("templates/weight.html", context)
	bit := []byte(str)
	w.Write(bit)
}
*/

func getMedication(w http.ResponseWriter, r *http.Request) {
	medications := []Medication{}
	db.Preload("Animals").Preload("Medicines").Find(&medications, Medication{})

	animals := []Animal{}
	db.Find(&animals, &Animal{})

	medicines := []Medicine{}
	db.Find(&medicines, &Medicine{})

	context := map[string]interface{}{
		"animals": animals,
		"medicines": medicines,
		"medications": medications,
	}

	str, _ := mustache.RenderFile("templates/medication.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postMedication(w http.ResponseWriter, r *http.Request) {	
	medication := Medication{}
	medication.Description = r.PostFormValue("Description")

	date, _ := time.Parse("2006-01-02", r.PostFormValue("Date"))
	medication.Date = mysql.NullTime{Time: date, Valid: true}

	r.ParseForm()
	for _, idAnimals := range r.Form["Animal"] {
		// element is the element from someSlice for where we are
		animal := Animal{}
		id, _ := strconv.Atoi(idAnimals)
		db.Find(&animal, id)
		medication.Animals = append(medication.Animals, animal)
	}

	r.ParseForm()
	for _, idMedicines := range r.Form["Medicine"] {
		// element is the element from someSlice for where we are
		medicine := Medicine{}
		id, _ := strconv.Atoi(idMedicines)
		db.Find(&medicine, id)
		medication.Medicines = append(medication.Medicines, medicine)
	}

	db.Save(&medication)

	http.Redirect(w, r, "/medication", http.StatusMovedPermanently)
}