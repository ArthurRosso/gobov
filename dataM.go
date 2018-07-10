package main

func DataM() {
	typeMedicine1 := TypeMedicine{}
	typeMedicine1.Type = "Via Oral"
	typeMedicine1.Description = "Caracterizada pela ingestão pela boca. Pode exercer efeitos locais no trato gastrointestinal ou atingir sangue e linfa provocando efeitos sistêmicos, após ser absorvido na mucosa gastrointestinal."
	db.Save(&typeMedicine1)

	typeMedicine2 := TypeMedicine{}
	typeMedicine2.Type = "Via Intradérmica"
	typeMedicine2.Description = "A injeção intradérmica consiste na aplicação de solução na derme (área localizada entre a derme e o tecido subcutâneo."
	db.Save(&typeMedicine2)

	typeMedicine3 := TypeMedicine{}
	typeMedicine3.Type = "Via Subcutânea"
	typeMedicine3.Description = "A injeção subcutânea consiste na aplicação de solução na região subcutânea, isto é, na hipoderme (tecido adiposo abaixo da pele)."
	db.Save(&typeMedicine3)

	typeMedicine4 := TypeMedicine{}
	typeMedicine4.Type = "Via Intramuscular"
	typeMedicine4.Description = "Consiste na aplicação de solução no tecido muscular."
	db.Save(&typeMedicine4)

	typeMedicine5 := TypeMedicine{}
	typeMedicine5.Type = "Via Pour-on"
	typeMedicine5.Description = "Consiste na aplicação sobre a linha média superior dos animais (espaço compreendido entre a cernelha e a inserção da cauda)."
	db.Save(&typeMedicine5)

	typeMedicine6 := TypeMedicine{}
	typeMedicine6.Type = "Via Spray"
	typeMedicine6.Description = "Consiste na aplicação pressionan a válvula do tubo e direcionan o jato para a região a ser tratada, pulverizando por alguns segundos."
	db.Save(&typeMedicine6)
}
