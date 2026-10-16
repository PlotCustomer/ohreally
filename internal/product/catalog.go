package product

// SeedCatalog returns the initial Product catalog of the Webshop. The catalog
// covers both BusinessObject specializations of the Product DataObject:
// Fruit & Groenten and Pasta.
func SeedCatalog() []Product {
	return []Product{
		{ID: "fruit-appel", Name: "Appels Elstar", Category: CategoryFruitGroenten, Description: "Knapperige Hollandse appels, per kilo.", PriceCents: 249, Currency: "EUR", Unit: "kg", Stock: 120},
		{ID: "fruit-banaan", Name: "Bananen", Category: CategoryFruitGroenten, Description: "Rijpe bananen, per tros.", PriceCents: 199, Currency: "EUR", Unit: "tros", Stock: 80},
		{ID: "fruit-tomaat", Name: "Tomaten", Category: CategoryFruitGroenten, Description: "Verse trostomaten van eigen bodem.", PriceCents: 329, Currency: "EUR", Unit: "kg", Stock: 60},
		{ID: "groente-aardappel", Name: "Aardappelen", Category: CategoryFruitGroenten, Description: "Kruimige aardappelen, ideaal voor het avondeten.", PriceCents: 179, Currency: "EUR", Unit: "kg", Stock: 200},
		{ID: "groente-wortel", Name: "Wortels", Category: CategoryFruitGroenten, Description: "Zoete winterwortels, per kilo.", PriceCents: 149, Currency: "EUR", Unit: "kg", Stock: 150},
		{ID: "pasta-spaghetti", Name: "Spaghetti", Category: CategoryPasta, Description: "Italiaanse spaghetti van durumtarwe, 500 gram.", PriceCents: 189, Currency: "EUR", Unit: "pak", Stock: 300},
		{ID: "pasta-penne", Name: "Penne", Category: CategoryPasta, Description: "Korte buisjes, perfect voor ovenschotels.", PriceCents: 179, Currency: "EUR", Unit: "pak", Stock: 250},
		{ID: "pasta-fusilli", Name: "Fusilli", Category: CategoryPasta, Description: "Gedraaide pasta die sauzen goed vasthoudt.", PriceCents: 199, Currency: "EUR", Unit: "pak", Stock: 220},
		{ID: "pasta-tagliatelle", Name: "Tagliatelle", Category: CategoryPasta, Description: "Lintpasta van verse eieren, 250 gram.", PriceCents: 269, Currency: "EUR", Unit: "pak", Stock: 140},
	}
}
