package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
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
	db, _ = gorm.Open("postgres", "postgres://qeympnnenynhnw:de0e9713a10d0da78775865f2196fd430f10dcd27b78e4c5c505c44e4d9ba339@ec2-107-20-211-10.compute-1.amazonaws.com:5432/dc9qbu79ln5lor")
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

	db.Table("type_medicines").Count(&countTypeMedicnes)
	if countTypeMedicnes == 0 {
		DataTypeMedicines()
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	logado := r.PathPrefix("/").Subrouter()
	logado.Use(loggingMiddleware)

	logado.HandleFunc("/", getIndex)

	logado.HandleFunc("/pic/{idAnimal}", getPic)
	logado.HandleFunc("/picMedicine/{idMedicine}", getMedicinePic)

	logado.HandleFunc("/animal", getAnimal)
	logado.HandleFunc("/newAnimal", postAnimal)
	logado.HandleFunc("/delAnimal/{ID}", delAnimal)
	logado.HandleFunc("/editAnimal/{ID}", editAnimal)
	logado.HandleFunc("/renewAnimal/{ID}", repostAnimal)
	logado.HandleFunc("/relatorioAnimal/{idAnimal}", relAnimal)
	logado.HandleFunc("/listaAnimal", getAllAnimals)

	logado.HandleFunc("/weight/{idAnimal}", getWeight)
	logado.HandleFunc("/newWeight/{idAnimal}", postWeight)
	logado.HandleFunc("/delWeight/{idWeight}", delWeight)

	logado.HandleFunc("/medicine", getMedicine)
	logado.HandleFunc("/newMedicine", postMedicine)
	logado.HandleFunc("/delMedicine/{ID}", delMedicine)
	logado.HandleFunc("/editMedicine/{ID}", editMedicine)
	logado.HandleFunc("/renewMedicine/{ID}", repostMedicine)
	logado.HandleFunc("/listaMedicine", getAllMedicines)

	logado.HandleFunc("/profile/{idAnimal}", getProfile)

	logado.HandleFunc("/medication", getMedication)
	logado.HandleFunc("/newMedication", postMedication)
	logado.HandleFunc("/delMedication/{ID}", delMedication)
	logado.HandleFunc("/listaMedication", getAllMedications)

	r.HandleFunc("/register", register)
	r.HandleFunc("/checkRegister", checkRegister)
	r.HandleFunc("/auth", auth)
	r.HandleFunc("/login", login)
	logado.HandleFunc("/logout", logout)

	logado.HandleFunc("/fazenda.xlsx", exportExcel)

	logado.HandleFunc("/preferences", preferences)
	logado.HandleFunc("/breed", getBreed)
	logado.HandleFunc("/newBreed", postBreed)
	logado.HandleFunc("/delBreed/{idBreed}", delBreed)

	logado.HandleFunc("/purp", getPurp)
	logado.HandleFunc("/newPurp", postPurp)
	logado.HandleFunc("/delPurp/{idPurp}", delPurp)

	logado.HandleFunc("/typea", getTypeA)
	logado.HandleFunc("/newTypeA", postTypeA)
	logado.HandleFunc("/delTypeA/{idType}", delTypeA)

	logado.HandleFunc("/typem", getTypeM)
	logado.HandleFunc("/newTypeM", postTypeM)
	logado.HandleFunc("/delTypeM/{idType}", delTypeM)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, r)
}

func exportExcel(w http.ResponseWriter, r *http.Request) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	animals := []Animal{}
	weights := []Weight{}
	medicines := []Medicine{}
	medications := []Medication{}
	ctx := GetContext(w, r)
	db.Where("user_id = ?", ctx.User.ID).Preload("Animals").Preload("Medicines").Find(&medications, Medication{})
	db.Where("user_id = ?", ctx.User.ID).Preload("Type").Find(&medicines, Medicine{})
	db.Where("user_id = ?", ctx.User.ID).Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Find(&animals, Animal{})

	file = xlsx.NewFile()

	sheet, _ = file.AddSheet("Animais")

	row = sheet.AddRow()
	name := row.AddCell()
	name.Value = "Identificador"
	birth := row.AddCell()
	birth.Value = "Data nasc"
	we := row.AddCell()
	we.Value = "Peso atual"
	ti := row.AddCell()
	ti.Value = "Tipo"
	br := row.AddCell()
	br.Value = "Raça"
	prop := row.AddCell()
	prop.Value = "Propósito"
	mom := row.AddCell()
	mom.Value = "Mãe"
	dad := row.AddCell()
	dad.Value = "Pai"

	fmt.Println("Cabeçalho animal não nulo")

	for _, animal := range animals {
		// element is the element from someSlice for where we are
		row = sheet.AddRow()
		name := row.AddCell()
		name.Value = animal.Name
		birth := row.AddCell()
		birth.Value = animal.BirthFmt()
		we := row.AddCell()
		we.Value = animal.WeightFmt()
		ti := row.AddCell()
		ti.Value = animal.Type.Type
		br := row.AddCell()
		br.Value = animal.Breed.Breed
		prop := row.AddCell()
		prop.Value = animal.PurposesFmt()
		mom := row.AddCell()
		mom.Value = animal.MotherFmt()
		dad := row.AddCell()
		dad.Value = animal.FatherFmt()
	}

	fmt.Println("Animais não nulos")

	for _, animal := range animals {
		sheet, _ = file.AddSheet("Peso animal " + animal.Name)
		row = sheet.AddRow()
		wei := row.AddCell()
		wei.Value = "Peso"
		desc := row.AddCell()
		desc.Value = "Descrição"
		data := row.AddCell()
		data.Value = "Data"

		fmt.Println("Cabeçalho peso animal não nulo")

		db.Where("animal_id = ?", animal.ID).Find(&weights)
		for _, weight := range weights {
			row = sheet.AddRow()
			wei := row.AddCell()
			wei.Value = fmt.Sprintf("%f3", weight.Weight)
			desc := row.AddCell()
			desc.Value = weight.Description
			date := row.AddCell()
			date.Value = weight.DateFmt()
		}
	}

	fmt.Println("Peso animal não nulo")

	sheet, _ = file.AddSheet("Remédios")
	row = sheet.AddRow()
	name = row.AddCell()
	name.Value = "Nome"
	exp := row.AddCell()
	exp.Value = "Validade"
	desc := row.AddCell()
	desc.Value = "Descrição"
	ti = row.AddCell()
	ti.Value = "Tipo"

	fmt.Println("Cabeçalho remédio não nulo")

	for _, medicine := range medicines {
		// element is the element from someSlice for where we are
		row = sheet.AddRow()
		name := row.AddCell()
		name.Value = medicine.Name
		exp := row.AddCell()
		exp.Value = medicine.ExpirationFmt()
		desc := row.AddCell()
		desc.Value = medicine.Description
		ti := row.AddCell()
		ti.Value = medicine.Type.Type
	}

	fmt.Println("Remédios não nulos")

	sheet, _ = file.AddSheet("Medicações")
	row = sheet.AddRow()
	desc = row.AddCell()
	desc.Value = "Descrição"
	date := row.AddCell()
	date.Value = "Data"
	an := row.AddCell()
	an.Value = "Animais"
	med := row.AddCell()
	med.Value = "Remédios"

	fmt.Println("Cabeçalho medicação não nulo")

	for _, medication := range medications {
		row = sheet.AddRow()
		desc := row.AddCell()
		desc.Value = medication.Description
		date := row.AddCell()
		date.Value = medication.DateFmt()
		an = row.AddCell()
		an.Value = medication.AnimalsFmt()
		med = row.AddCell()
		med.Value = medication.MedicinesFmt()
	}

	fmt.Println("Medicações não nulas")

	file.Write(w)

	http.Redirect(w, r, "/", http.StatusFound)
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
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	password = string(bytes)
	db.Where("username = ? AND password = ?", username, password).First(&user, User{})

	if user.Username != "" {
		http.Redirect(w, r, "/register", http.StatusFound)
	} else {
		user.Username = username
		user.Password = password
		user.Email = r.PostFormValue("Email")
		db.Save(&user)

		// Set user as authenticated
		db.First(&user)

		ctx := GetContext(w, r)
		ctx.Session.Values["User.ID"] = user.ID
		ctx.Close()
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	err := ctx.GetFlashes()
	ctx.Close()
	context := map[string]interface{}{
		"err": err,
	}

	str, _ := mustache.RenderFile("templates/login.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func auth(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)

	// Authentication goes here
	user := User{}
	username := r.PostFormValue("Username")
	password := r.PostFormValue("Password")
	db.Where("username = ?", username).First(&user, User{})

	if user.ID == 0 || nil != bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) {
		ctx.AddFlash("Nome de usuário ou senha incorreto(s)")
		ctx.Close()
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
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
	db.Where("user_id = ?", ctx.User.ID).Order("id desc").Find(&histories, History{})

	context := map[string]interface{}{
		"histories": histories,
		"user":      ctx.User,
		"counta":    countAnimals,
		"countm":    countMedicines,
		"countmed":  countMedications,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/index.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/animal.html", context)
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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/listAnimal.html", context)
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
	animal.Birthday = birth

	weight := Weight{
		Description: "Primeira pesagem",
		Date:        birth,
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
	history.Description = "Cadastro de animal realizado: " + name
	history.User = ctx.User
	history.Animals = []*Animal{&animal}
	history.Date = time.Now()

	db.Save(&history)

	db.Save(&animal)

	if files == nil {
		pic := Picture{Main: true, AnimalID: animal.ID}
		f, _ := ioutil.ReadFile("static/cow-and-moon.jpg")
		pic.Picture = f
		db.Save(&pic)
	}

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

	http.Redirect(w, r, "/profile/"+strconv.Itoa(animal.ID), http.StatusFound)
}

func editAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	animal := Animal{ID: id}
	db.Where("user_id = ?", ctx.User.ID).Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Find(&animal, Animal{})

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
		"animal":   animal,
		"types":    types,
		"breeds":   breeds,
		"purposes": purposes,
		"animals":  animals,
		"mothers":  mothers,
		"fathers":  fathers,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/editAnimal.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func repostAnimal(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(w, r)

	mvar := mux.Vars(r)
	id, _ := strconv.Atoi(mvar["ID"])
	animal := Animal{ID: id}
	db.Find(&animal, &Animal{})

	db.Preload("Medications").Preload("Purposes").First(&animal)

	db.Exec("DELETE FROM weights WHERE animal_id=?", id)
	db.Exec("DELETE FROM animal_purpose WHERE animal_id=?", id)
	db.Exec("DELETE FROM medication_animal WHERE animal_id=?", id)

	db.Where("animal_id = ?", id).Delete(&Picture{})

	db.Delete(&animal)

	animal = NewAnimal()
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
	animal.Birthday = birth

	weight := Weight{
		Description: "Primeira pesagem",
		Date:        birth,
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

	history := History{}
	history.Description = "Animal editado: " + name
	history.User = ctx.User
	history.Animals = []*Animal{&animal}
	history.Date = time.Now()

	animal.User = ctx.User

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

	http.Redirect(w, r, "/listaAnimal", http.StatusFound)
}

func delAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)

	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	animal := Animal{ID: id}
	db.Find(&animal, &Animal{})

	history := History{}
	history.Description = "Exclusão realizada do animal: " + animal.Name
	history.User = ctx.User
	history.Animals = []*Animal{&animal}
	history.Date = time.Now()
	db.Save(&history)

	db.Preload("Medications").Preload("Purposes").First(&animal)

	db.Exec("DELETE FROM weights WHERE animal_id=?", id)
	db.Exec("DELETE FROM animal_purpose WHERE animal_id=?", id)
	db.Exec("DELETE FROM medication_animal WHERE animal_id=?", id)

	db.Where("animal_id = ?", id).Delete(&Picture{})
	db.Delete(&animal)

	http.Redirect(w, r, "/listaAnimal", http.StatusFound)
}

func delMedicine(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	medicine := Medicine{ID: id}
	db.Preload("Medications").First(&medicine)

	ctx := GetContext(w, r)

	history := History{}
	history.Description = "Exclusão realizada do remédio: " + medicine.Name
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Exec("DELETE FROM medication_medicine WHERE medicine_id=?", id)
	db.Where("ID = ?", id).Delete(&Medicine{})

	http.Redirect(w, r, "/listaMedicine", http.StatusFound)
}

func getMedicine(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)

	medicines := []Medicine{}
	db.Where("user_id = ?", ctx.User.ID).Preload("Type").Find(&medicines, Medicine{})

	types := []TypeMedicine{}
	db.Find(&types, &TypeMedicine{})

	context := map[string]interface{}{
		"types":     types,
		"medicines": medicines,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/medicine.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func editMedicine(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	medicine := Medicine{ID: id}
	db.Where("user_id = ?", ctx.User.ID).Preload("Type").First(&medicine, Medicine{})

	types := []TypeMedicine{}
	db.Find(&types, &TypeMedicine{})

	context := map[string]interface{}{
		"types":    types,
		"medicine": medicine,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/editMedicine.html", context)
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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/listMedicine.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postMedicine(w http.ResponseWriter, r *http.Request) {
	medicine := NewMedicine()
	name := r.PostFormValue("Name")
	medicine.Name = name

	expiration, _ := time.Parse("2006-01-02", r.PostFormValue("Expiration"))
	medicine.Expiration = expiration

	medicine.Description = r.PostFormValue("Description")

	typeM := TypeMedicine{}
	idType, _ := strconv.Atoi(r.PostFormValue("Type"))
	db.Find(&typeM, idType)
	medicine.Type = &typeM
	db.First(&medicine.Type, idType)

	ctx := GetContext(w, r)
	medicine.User = ctx.User

	history := History{}
	history.Description = "Cadastro de remédio realizado: " + name
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	r.ParseMultipartForm(0)
	f := r.MultipartForm
	file := f.File["Picture"]
	if file != nil {
		arquivo, _ := file[0].Open()
		medicine.Picture, _ = ioutil.ReadAll(arquivo)
		arquivo.Close()
	}
	db.Save(&medicine)

	http.Redirect(w, r, "/listaMedicine", http.StatusFound)
}

func repostMedicine(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	medicine := Medicine{ID: id}
	db.Preload("Medications").First(&medicine)

	db.Exec("DELETE FROM medication_medicine WHERE medicine_id=?", id)
	db.Where("ID = ?", id).Delete(&Medicine{})

	medicine = NewMedicine()
	name := r.PostFormValue("Name")
	medicine.Name = name

	expiration, _ := time.Parse("2006-01-02", r.PostFormValue("Expiration"))
	medicine.Expiration = expiration

	medicine.Description = r.PostFormValue("Description")

	typeM := TypeMedicine{}
	idType, _ := strconv.Atoi(r.PostFormValue("Type"))
	db.Find(&typeM, idType)
	medicine.Type = &typeM
	db.First(&medicine.Type, idType)

	r.ParseMultipartForm(0)
	f := r.MultipartForm
	if f == nil {
		fmt.Println("_o no formulário")
	}
	file := f.File["Picture"]
	arquivo, _ := file[0].Open()
	medicine.Picture, _ = ioutil.ReadAll(arquivo)
	defer arquivo.Close()

	medicine.User = ctx.User

	db.Save(&medicine)

	http.Redirect(w, r, "/listaMedicine", http.StatusFound)
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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/medication.html", context)
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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/listMedication.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postMedication(w http.ResponseWriter, r *http.Request) {
	medication := Medication{}
	desc := r.PostFormValue("Description")
	medication.Description = desc

	date, _ := time.Parse("2006-01-02", r.PostFormValue("Date"))
	medication.Date = date

	r.ParseForm()
	for _, idAnimals := range r.Form["Animal"] {
		animal := Animal{}
		id, _ := strconv.Atoi(idAnimals)
		db.Find(&animal, id)
		medication.Animals = append(medication.Animals, &animal)
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
	history.Animals = medication.Animals
	history.User = ctx.User
	history.Medication = &medication
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&medication)

	http.Redirect(w, r, "/listaMedication", http.StatusFound)
}

func delMedication(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)

	m := mux.Vars(r)
	id, _ := strconv.Atoi(m["ID"])
	medication := Medication{ID: id}
	db.Find(&medication, &Medication{})

	history := History{}
	history.Description = "Exclusão de medicação realizada: " + medication.Description
	history.User = ctx.User
	history.Animals = medication.Animals
	t := time.Now()
	history.Date = t
	db.Save(&history)

	db.Preload("Animals").Preload("Medicines").First(&medication)

	db.Exec("DELETE FROM medication_animal WHERE medication_id=?", id)
	db.Exec("DELETE FROM medication_medicine WHERE medication_id=?", id)

	db.Delete(&medication)

	http.Redirect(w, r, "/listaMedication", http.StatusFound)
}

func getPic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	picture := Picture{}
	animal := Animal{}
	db.First(&animal, idAnimal)
	db.First(&picture, idAnimal)
	if len(picture.Picture) > 0 {
		w.Write(picture.Picture)
	}
}

func getMedicinePic(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	vars := mux.Vars(r)
	medicine := Medicine{}
	idMedicine, _ := strconv.Atoi(vars["idMedicine"])
	db.First(&medicine, idMedicine)
	if medicine.UserID != ctx.User.ID {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if len(medicine.Picture) > 0 {
		w.Write(medicine.Picture)
	}
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{ID: idAnimal}
	db.Preload("Medications").Preload("Weights").Preload("Type").Preload("Breed").Preload("Purposes").Preload("Father").Preload("Mother").Preload("Pictures").First(&animal, idAnimal)

	if animal.UserID != ctx.User.ID {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	histories := []History{}
	db.Model(&animal).Order("id desc").Related(&histories, "Histories")

	context := map[string]interface{}{
		"histories": histories,
		"animal":    animal,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/profile.html", context)
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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/charts.html", context)
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

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/weight.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postWeight(w http.ResponseWriter, r *http.Request) {
	weight := Weight{}
	peso, _ := strconv.ParseFloat(r.PostFormValue("Weight"), 32)
	weight.Weight = float32(peso)
	desc := r.PostFormValue("Description")
	weight.Description = desc

	date, _ := time.Parse("2006-01-02", r.PostFormValue("Date"))
	weight.Date = date

	vars := mux.Vars(r)
	idAnimal, _ := strconv.Atoi(vars["idAnimal"])
	animal := Animal{}
	db.First(&animal, idAnimal)

	weight.Animal = &animal

	ctx := GetContext(w, r)
	history := History{}
	history.Description = "Pesagem realizada: " + desc + " do animal " + animal.Name
	history.User = ctx.User
	history.Animals = []*Animal{&animal}
	history.Date = time.Now()
	db.Save(&history)

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

	ctx := GetContext(w, r)
	history := History{}
	animal := Animal{ID: weight.AnimalID}
	db.First(&animal, weight.AnimalID)
	history.Description = "Exclusão de pesagem realizada: " + weight.Description + " do animal " + animal.Name
	history.User = ctx.User
	history.Animals = []*Animal{&animal}
	history.Date = time.Now()
	db.Save(&history)

	db.Delete(&weight)

	url := "/weight/" + id

	http.Redirect(w, r, url, http.StatusFound)
}

func preferences(w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/preferences.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func getBreed(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	breeds := []Breed{}
	db.Where("user_id = 0 OR user_id = ?", ctx.User.ID).Find(&breeds, &Breed{})

	context := map[string]interface{}{
		"breeds": breeds,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/breed.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postBreed(w http.ResponseWriter, r *http.Request) {
	breed := Breed{}
	breed.Breed = r.PostFormValue("Breed")
	breed.Description = r.PostFormValue("Description")

	ctx := GetContext(w, r)
	breed.User = ctx.User
	history := History{}
	history.Description = "Raça cadastrada: " + breed.Breed
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&breed)

	http.Redirect(w, r, "/breed", http.StatusFound)
}

func delBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idBreed, _ := strconv.Atoi(vars["idBreed"])
	breed := Breed{}
	db.Find(&breed, idBreed)

	ctx := GetContext(w, r)
	history := History{}
	history.Description = "Exclusão de raça realizada: " + breed.Breed
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Delete(&breed)

	http.Redirect(w, r, "/breed", http.StatusFound)
}

func getPurp(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	purps := []Purpose{}
	db.Where("user_id = 0 OR user_id = ?", ctx.User.ID).Find(&purps, &Purpose{})

	context := map[string]interface{}{
		"purps": purps,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/purp.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postPurp(w http.ResponseWriter, r *http.Request) {
	purp := Purpose{}
	purp.Purpose = r.PostFormValue("Purpose")
	purp.Description = r.PostFormValue("Description")

	ctx := GetContext(w, r)
	purp.User = ctx.User
	history := History{}
	history.Description = "Finalidade cadastrada: " + purp.Purpose
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&purp)

	http.Redirect(w, r, "/purp", http.StatusFound)
}

func delPurp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPurp, _ := strconv.Atoi(vars["idPurp"])
	purp := Purpose{}
	db.Find(&purp, idPurp)

	ctx := GetContext(w, r)
	history := History{}
	history.Description = "Exclusão de finalidade realizada: " + purp.Purpose
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Delete(&purp)

	http.Redirect(w, r, "/purp", http.StatusFound)
}

func getTypeA(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	types := []TypeAnimal{}
	db.Where("user_id = 0 OR user_id = ?", ctx.User.ID).Find(&types, &TypeAnimal{})

	context := map[string]interface{}{
		"types": types,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/typea.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postTypeA(w http.ResponseWriter, r *http.Request) {
	typea := TypeAnimal{}
	typea.Type = r.PostFormValue("Type")
	typea.Description = r.PostFormValue("Description")

	ctx := GetContext(w, r)
	typea.User = ctx.User
	history := History{}
	history.Description = "Tipo de animal cadastrado: " + typea.Type
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&typea)

	http.Redirect(w, r, "/typea", http.StatusFound)
}

func delTypeA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idType, _ := strconv.Atoi(vars["idType"])
	typea := TypeAnimal{}
	db.Find(&typea, idType)

	ctx := GetContext(w, r)
	history := History{}
	history.Description = "Exclusão de tipo de animal realizado: " + typea.Type
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Delete(&typea)

	http.Redirect(w, r, "/typea", http.StatusFound)
}

func getTypeM(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(w, r)
	types := []TypeMedicine{}
	db.Where("user_id = 0 OR user_id = ?", ctx.User.ID).Find(&types, &TypeMedicine{})

	context := map[string]interface{}{
		"types": types,
	}

	str, _ := mustache.RenderFileInLayout("templates/navbar.template.html", "templates/typem.html", context)
	bit := []byte(str)
	w.Write(bit)
}

func postTypeM(w http.ResponseWriter, r *http.Request) {
	typem := TypeMedicine{}
	typem.Type = r.PostFormValue("Type")
	typem.Description = r.PostFormValue("Description")

	ctx := GetContext(w, r)
	typem.User = ctx.User
	history := History{}
	history.Description = "Tipo de remédio cadastrado: " + typem.Type
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Save(&typem)

	http.Redirect(w, r, "/typem", http.StatusFound)
}

func delTypeM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idType, _ := strconv.Atoi(vars["idType"])
	typem := TypeMedicine{}
	db.Find(&typem, idType)

	ctx := GetContext(w, r)
	history := History{}
	history.Description = "Exclusão de tipo de remédio realizado: " + typem.Type
	history.User = ctx.User
	history.Date = time.Now()
	db.Save(&history)

	db.Delete(&typem)

	http.Redirect(w, r, "/typem", http.StatusFound)
}
