package main

func Data() {
	typeAnimal1 := TypeAnimal{}
	typeAnimal1.Type = "Touro"
	typeAnimal1.Description = "Macho bovino não castrado"
	db.Save(&typeAnimal1)

	typeAnimal2 := TypeAnimal{}
	typeAnimal2.Type = "Vaca"
	typeAnimal2.Description = "Fêmea bovina"
	db.Save(&typeAnimal2)

	typeAnimal3 := TypeAnimal{}
	typeAnimal3.Type = "Boi"
	typeAnimal3.Description = "Macho bovino castrado"
	db.Save(&typeAnimal3)
	breed1 := Breed{}
	breed1.Breed = "Angus"
	breed1.Description = "O Aberdeen Angus se destaca entre as raças taurinas por reunir um maior número de características positivas que lhe asseguram um excelente resultado econômico como gado de corte. O conjunto de suas características a tornam uma raça completa."
	db.Save(&breed1)

	purpose1 := Purpose{}
	purpose1.Purpose = "Genética"
	purpose1.Description = "Animal destinado a gerar descendentes.'"
	db.Save(&purpose1)

	purpose2 := Purpose{}
	purpose2.Purpose = "Leite"
	purpose2.Description = "Animal destinado a produção de leite."
	db.Save(&purpose2)

	purpose3 := Purpose{}
	purpose3.Purpose = "Engorda"
	purpose3.Description = "Animal destinado a engorda para produção de carne."
	db.Save(&purpose3)
}
