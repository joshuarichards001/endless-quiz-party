package main

var categories = map[string][]string{
	"Science": {
		"Physics", "Chemistry", "Biology", "Geology", "Environmental Science", "Human Body", "Inventions", "Mathematics", "Famous Scientists",
	},
	"Geography": {
		"World Capitals", "Landmarks", "Physical Geography", "Flags of the World", "Countries and Cultures", "Oceans and Seas", "Mountains and Volcanoes", "Rivers and Lakes", "Continents", "Deserts",
	},
	"History": {
		"Ancient Civilizations", "Medieval History", "Modern History", "War History", "Historical Figures", "World Wars", "Revolutions", "Explorers", "Inventions in History", "Political Leaders",
	},
	"Literature": {
		"Classic Literature", "Modern Literature", "Poets and Authors", "Book Characters", "Famous Quotations", "Children's Books", "Literary Awards", "World Mythology", "Shakespeare", "Literary Genres",
	},
	"Arts & Entertainment": {
		"Movies", "Television Shows", "Music", "Visual Arts", "Theater", "Famous Paintings", "Classical Music", "Pop Culture", "Oscars and Awards", "Cartoons and Animation",
	},
	"Sports": {
		"Football (Soccer)", "American Football", "Baseball", "Basketball", "Olympics", "Extreme Sports", "Tennis", "Golf", "Cricket", "Winter Sports",
	},
	"Technology": {
		"Computers and Internet", "Gadgets and Devices", "AI and Robotics", "Space Technology", "Video Games", "Tech Companies", "Mobile Phones", "Social Media", "Cybersecurity",
	},
	"General Knowledge": {
		"Current Affairs", "Popular Culture", "Food and Drink", "Etiquette and Manners", "World Records", "Famous Brands", "Holidays and Traditions", "Language and Words", "Measurement and Units", "Famous Quotes",
	},
	"Nature & Environment": {
		"Wildlife", "Plants and Trees", "Natural Disasters", "Habitats", "Conservation", "Weather and Climate", "Endangered Species", "Oceans", "National Parks", "Ecosystems",
	},
	"Business & Economics": {
		"Famous Companies", "Currencies", "Economists", "Stock Markets", "Business Leaders", "Brands and Logos", "Economic Terms", "Trade and Commerce", "Inventions in Business", "Advertising",
	},
	"Society & Culture": {
		"World Religions", "Philosophy", "Languages", "Customs and Traditions", "Fashion", "Famous Landmarks", "Festivals", "Social Movements", "Human Rights", "Famous Speeches",
	},
	"Transportation": {
		"Cars", "Trains", "Aviation", "Ships and Boats", "Public Transport", "Space Travel", "Famous Explorers", "Bridges and Tunnels", "Roads and Highways", "Inventions in Transport",
	},
	"Health & Medicine": {
		"Human Body", "Diseases", "Medical Discoveries", "Nutrition", "Famous Doctors", "Mental Health", "First Aid", "Medical Technology", "Vaccines", "Public Health",
	},
	"Language & Words": {
		"English Language", "Word Origins", "Idioms and Phrases", "Foreign Languages", "Spelling", "Grammar", "Famous Authors", "Proverbs", "Synonyms and Antonyms",
	},
	"Space & Astronomy": {
		"Planets", "Stars", "Galaxies", "Space Missions", "Astronauts", "Telescopes", "The Moon", "The Sun", "Space Exploration", "Comets and Asteroids",
	},
}
