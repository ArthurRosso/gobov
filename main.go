package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	db                *gorm.DB
	countAnimals      int
	countMedicines    int
	countMedications  int
	countBreeds       int
	countPurposes     int
	countTypeAnimals  int
	countTypeMedicnes int
)

func main() {
	db, _ = gorm.Open("mysql", "goboi:goboi@/goboi")
	defer db.Close()

	db.AutoMigrate(&Animal{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Weight{})
	db.AutoMigrate(&TypeAnimal{})
	db.AutoMigrate(&Breed{})
	db.AutoMigrate(&Purpose{})
	db.AutoMigrate(&Medicine{})
	db.AutoMigrate(&TypeMedicine{})
	db.AutoMigrate(&Picture{})
	db.AutoMigrate(&Medication{})
	db.AutoMigrate(&History{})

	r := mux.NewRouter()

	logado := r.PathPrefix("/").Subrouter()
	logado.Use(loggingMiddleware)

	logado.HandleFunc("/", getIndex)

	logado.HandleFunc("/pic/{idAnimal}", getPic)
	logado.HandleFunc("/picMedicine/{idMedicine}", getMedicinePic)

	logado.HandleFunc("/animal", getAnimal)
	logado.HandleFunc("/newAnimal", postAnimal)
	logado.HandleFunc("/delAnimal/{ID}", delAnimal)
	//logado.HandleFunc("/editAnimal/{ID}", editAnimal)
	logado.HandleFunc("/relatorioAnimal/{idAnimal}", relAnimal)
	logado.HandleFunc("/listaAnimal", getAllAnimals)

	logado.HandleFunc("/weight/{idAnimal}", getWeight)
	logado.HandleFunc("/newWeight/{idAnimal}", postWeight)
	logado.HandleFunc("/delWeight/{idWeight}", delWeight)

	logado.HandleFunc("/medicine", getMedicine)
	logado.HandleFunc("/newMedicine", postMedicine)
	logado.HandleFunc("/delMedicine/{ID}", delMedicine)
	logado.HandleFunc("/listaMedicine", getAllMedicines)

	logado.HandleFunc("/profile/{idAnimal}", getProfile)

	logado.HandleFunc("/medication", getMedication)
	logado.HandleFunc("/newMedication", postMedication)
	logado.HandleFunc("/listaMedication", getAllMedications)

	r.HandleFunc("/register", register)
	r.HandleFunc("/checkRegister", checkRegister)
	r.HandleFunc("/auth", auth)
	r.HandleFunc("/login", login)
	logado.HandleFunc("/logout", logout)

	fmt.Println("Server listen and serve on port 8000")

	http.ListenAndServe(":8000", r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := GetContext(w, r)

		if ctx.User == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		}
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{}

	str, _ := mustache.RenderFile("templates/register.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func checkRegister(w http.ResponseWriter, r *http.Request) {

	// Authentication goes here
	user := User{}
	username := r.PostFormValue("Username")
	password := r.PostFormValue("Password")
	db.Where("username = ? AND password = ?", username, password).First(&user, User{})

	if user.Username != "" {
		http.Redirect(w, r, "/register", http.StatusFound)
	} else {
		user.Username = username
		user.Password = password
		user.Name = r.PostFormValue("Name")
		user.Email = r.PostFormValue("Email")
		db.Save(&user)

		// Set user as authenticated
		db.First(&user)

		ctx := GetContext(w, r)
		defer ctx.Close()
		ctx.Session.Values["User.ID"] = user.ID

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{}

	str, _ := mustache.RenderFile("templates/login.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func auth(w http.ResponseWriter, r *http.Request) {

	// Authentication goes here
	user := User{}
	username := r.PostFormValue("Username")
	password := r.PostFormValue("Password")
	db.Where("username = ? AND password = ?", username, password).First(&user, User{})

	if user.Username == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		ctx := GetContext(w, r)
		ctx.Session.Values["User.ID"] = user.ID
		ctx.Close()

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	ctx.Session.Values = map[interface{}]interface{}{}

	ctx.Close()

	http.Redirect(w, r, "/login", http.StatusFound)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)

	db.Where("user_id = ?", ctx.User.ID).Table("animals").Count(&countAnimals)
	db.Where("user_id = ?", ctx.User.ID).Table("medicines").Count(&countMedicines)
	db.Where("user_id = ?", ctx.User.ID).Table("medications").Count(&countMedications)

	histories := []History{}
	db.Where("user_id = ?", ctx.User.ID).Find(&histories, History{})

	context := map[string]interface{}{
		"histories": histories,
		"user":      ctx.User,
		"counta":    countAnimals,
		"countm":    countMedicines,
		"countmed":  countMedications,
	}

	str, _ := mustache.RenderFile("templates/index.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	db.Table("breeds").Count(&countBreeds)
	if countBreeds == 0 {
		DataBreeds()
	}
	db.Table("purposes").Count(&countPurposes)
	if countPurposes == 0 {
		DataPurposes()
	}
	db.Table("type_animals").Count(&countTypeAnimals)
	if countTypeAnimals == 0 {
		DataTypeAnimals()
	}

	animals := []Animal{}
	db.Where("user_id = ?", ctx.User.ID).Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Find(&animals, Animal{})

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

func getAllAnimals(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	animals := []Animal{}
	db.Where("user_id = ?", ctx.User.ID).Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Find(&animals, Animal{})

	context := map[string]interface{}{
		"animals": animals,
	}

	str, _ := mustache.RenderFile("templates/listAnimal.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postAnimal(w http.ResponseWriter, r *http.Request) {
	animal := NewAnimal()
	name := r.PostFormValue("Name")
	animal.Name = name

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

	birth, _ := time.Parse("2006-01-02", r.PostFormValue("Birthday"))
	b := mysql.NullTime{Time: birth, Valid: true}
	animal.Birthday = b

	weight := Weight{
		Description: "Primeira pesagem",
		Date:        b,
	}
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

	ctx := GetContext(w, r)
	animal.User = ctx.User

	history := History{}
	history.Description = "Cadastro realizado: " + name
	history.User = ctx.User
	history.Animal = &animal
	history.Date = time.Now()

	db.Save(&history)

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

	http.Redirect(w, r, "/animal", http.StatusFound)
}

func delAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)

	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	animal := Animal{ID: id}

	history := History{}
	history.Description = "Exclusão realizada: " + animal.Name
	history.User = ctx.User
	history.Animal = &animal
	history.Date = time.Now()
	db.Save(&history)

	db.Preload("Medications").Preload("Purposes").First(&animal)

	db.Exec("DELETE FROM weights WHERE animal_id=?", id)
	db.Exec("DELETE FROM animal_purpose WHERE animal_id=?", id)
	db.Exec("DELETE FROM medication_animal WHERE animal_id=?", id)

	db.Where("animal_id = ?", id).Delete(&Picture{})
	db.Delete(&animal)

	http.Redirect(w, r, "/animal", http.StatusFound)
}

func delMedicine(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	medicine := Medicine{ID: id}
	db.Preload("Medications").First(&medicine)

	ctx := GetContext(w, r)

	history := History{}
	history.Description = "Exclusão realizada: " + medicine.Name
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Exec("DELETE FROM medication_medicine WHERE medicine_id=?", id)
	db.Where("ID = ?", id).Delete(&Medicine{})

	http.Redirect(w, r, "/medicine", http.StatusFound)
}

func getMedicine(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	db.Table("type_medicines").Count(&countTypeMedicnes)
	if countTypeMedicnes == 0 {
		DataTypeMedicines()
	}
	medicines := []Medicine{}
	db.Where("user_id = ?", ctx.User.ID).Preload("Type").Find(&medicines, Medicine{})

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

func getAllMedicines(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	medicines := []Medicine{}
	db.Where("user_id = ?", ctx.User.ID).Preload("Type").Find(&medicines, Medicine{})

	context := map[string]interface{}{
		"medicines": medicines,
	}

	str, _ := mustache.RenderFile("templates/listMedicine.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postMedicine(w http.ResponseWriter, r *http.Request) {
	medicine := NewMedicine()
	name := r.PostFormValue("Name")
	medicine.Name = name

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

	ctx := GetContext(w, r)
	medicine.User = ctx.User

	history := History{}
	history.Description = "Cadastro realizado: " + name
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&medicine)

	http.Redirect(w, r, "/medicine", http.StatusFound)
}

func getMedication(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	animals := []Animal{}
	db.Where("user_id = ?", ctx.User.ID).Find(&animals, &Animal{})

	medicines := []Medicine{}
	db.Where("user_id = ?", ctx.User.ID).Find(&medicines, &Medicine{})

	context := map[string]interface{}{
		"animals":   animals,
		"medicines": medicines,
	}

	str, _ := mustache.RenderFile("templates/medication.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func getAllMedications(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	medications := []Medication{}
	db.Where("user_id = ?", ctx.User.ID).Preload("Animals").Preload("Medicines").Find(&medications, Medication{})

	context := map[string]interface{}{
		"medications": medications,
	}

	str, _ := mustache.RenderFile("templates/listMedication.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postMedication(w http.ResponseWriter, r *http.Request) {
	medication := Medication{}
	desc := r.PostFormValue("Description")
	medication.Description = desc

	date, _ := time.Parse("2006-01-02", r.PostFormValue("Date"))
	medication.Date = mysql.NullTime{Time: date, Valid: true}

	r.ParseForm()
	for _, idAnimals := range r.Form["Animal"] {
		animal := Animal{}
		id, _ := strconv.Atoi(idAnimals)
		db.Find(&animal, id)
		medication.Animals = append(medication.Animals, animal)
	}

	r.ParseForm()
	for _, idMedicines := range r.Form["Medicine"] {
		medicine := Medicine{}
		id, _ := strconv.Atoi(idMedicines)
		db.Find(&medicine, id)
		medication.Medicines = append(medication.Medicines, medicine)
	}

	ctx := GetContext(w, r)
	medication.User = ctx.User

	history := History{}
	history.Description = "Medicação realizada: " + desc
	history.User = ctx.User
	history.Medication = &medication
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&medication)

	http.Redirect(w, r, "/medication", http.StatusFound)
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

func getMedicinePic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	medicine := Medicine{}
	idMedicine, _ := strconv.Atoi(vars["idMedicine"])
	db.First(&medicine, idMedicine)
	if len(medicine.Picture) > 0 {
		w.Write(medicine.Picture)
	}
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{ID: idAnimal}
	db.Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Preload("Father").Preload("Mother").Preload("Pictures").First(&animal, idAnimal)

	histories := []History{}

	medication := Medication{}
	db.Model(&animal).Related(&medication, "Medications")

	db.Where("animal_id = ?", animal.ID).Or("medication_id = ?", medication.ID).Find(&histories, History{})

	context := map[string]interface{}{
		"histories": histories,
		"animal":    animal,
	}

	str, _ := mustache.RenderFile("templates/profile.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func relAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{}
	db.First(&animal, idAnimal)

	weights := []Weight{}
	db.Where("animal_id = ?", idAnimal).Find(&weights)

	context := map[string]interface{}{
		"weights": weights,
		"animal":  animal,
	}

	str, _ := mustache.RenderFile("templates/charts.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func getWeight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{ID: idAnimal}

	weights := []Weight{}
	db.Where("animal_id = ?", idAnimal).Find(&weights)

	context := map[string]interface{}{
		"weights": weights,
		"animal":  animal,
	}

	str, _ := mustache.RenderFile("templates/weight.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postWeight(w http.ResponseWriter, r *http.Request) {
	weight := Weight{}
	peso, _ := strconv.ParseFloat(r.PostFormValue("Weight"), 32)
	weight.Weight = float32(peso)
	weight.Description = r.PostFormValue("Description")

	date, _ := time.Parse("2006-01-02", r.PostFormValue("Date"))
	weight.Date = mysql.NullTime{Time: date, Valid: true}

	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{}
	db.First(&animal, idAnimal)

	weight.Animal = &animal

	db.Save(&weight)

	url := "/weight/" + vars["idAnimal"]

	http.Redirect(w, r, url, http.StatusFound)
}

func delWeight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idWeight, _ := strconv.Atoi(vars["idWeight"])
	weight := Weight{}
	db.Find(&weight, idWeight)
	id := strconv.Itoa(weight.AnimalID)
	db.Delete(&weight)

	url := "/weight/" + id

	http.Redirect(w, r, url, http.StatusFound)
}
