package main

func DataM() {
	typeMedicine1 := TypeMedicine{}
	typeMedicine1.Type = "Via oral"
	typeMedicine1.Description = "Caracterizada pela ingestão pela boca. Pode exercer efeitos locais no trato gastrointestinal ou atingir sangue e linfa provocando efeitos sistêmicos, após ser absorvido na mucosa gastrointestinal."
	db.Save(&typeMedicine1)

	typeMedicine2 := TypeMedicine{}
	typeMedicine2.Type = "Via intradérmica"
	typeMedicine2.Description = "A injeção intradérmica consiste na aplicação de solução na derme (área localizada entre a derme e o tecido subcutâneo."
	db.Save(&typeMedicine2)

	typeMedicine3 := TypeMedicine{}
	typeMedicine3.Type = "Via subcutânea"
	typeMedicine3.Description = "A injeção subcutânea consiste na aplicação de solução na região subcutânea, isto é, na hipoderme (tecido adiposo abaixo da pele)."
	db.Save(&typeMedicine3)

	typeMedicine4 := TypeMedicine{}
	typeMedicine4.Type = "Via intramuscular"
	typeMedicine4.Description = "Consiste na aplicação de solução no tecido muscular."
	db.Save(&typeMedicine4)

}
