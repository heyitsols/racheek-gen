package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"github.com/gorilla/mux"
	"html/template"
	math_rand "math/rand"
	"net/http"
	"os"
)

type Character struct {
	Weapon, Background string
	Body, Mind, HP     int
	Items              []string
}

func roll() int {
	return math_rand.Intn(6-1) + 1
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

var items = map[int]string{
	11: "Alcohol",
	12: "Amulet",
	13: "Bell",
	14: "Blanket",
	15: "Bottle",
	16: "Bucket",
	21: "Candles",
	22: "Cards",
	23: "Chain",
	24: "Chalk",
	25: "Clay",
	26: "Compass",
	31: "Crowbar",
	32: "Dice",
	33: "Flint",
	34: "Hook",
	35: "Ink and paper",
	36: "Lockpicks",
	41: "Manacles",
	42: "Marbles",
	43: "Mirror",
	44: "Nails",
	45: "Oil",
	46: "Pans",
	51: "Perfume",
	52: "Pick",
	53: "Rope",
	54: "Sack",
	55: "Shovel",
	56: "Soap",
	61: "Spikes",
	62: "Telescope",
	63: "Torch",
	64: "Wax",
	65: "Whistle",
	66: "Wooden spoon",
}

var backgrounds = map[int]string{
	3:  "Baker",
	4:  "Barkeeper",
	5:  "Bounty Hunter",
	6:  "Herbalist",
	7:  "Investigator",
	8:  "Mercenary",
	9:  "Merchant",
	10: "Noble",
	11: "Pirate",
	12: "Priest",
	13: "Ranger",
	14: "Sailor",
	15: "Scholar",
	16: "Soldier",
	17: "Thief",
	18: "Traveler",
}

var weapons = map[int]string{
	1: "Unarmed (1)",
	2: "Basic (d6)",
	3: "Basic (d6)",
	4: "Basic (d6)",
	5: "Decent (2d6)",
	6: "Decent (2d6)",
}

var generate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var b [16]byte
	_, _ = crypto_rand.Read(b[:])
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	body := roll()
	mind := roll()
	hp := body + roll()
	itemCount := roll()
	weapon := weapons[roll()]
	background := backgrounds[roll()+roll()+roll()]

	var itemsGot []string

	for i := itemCount; i > 0; i-- {
		itemCode := (roll() * 10) + roll()
		if !contains(itemsGot, items[itemCode]) {
			itemsGot = append(itemsGot, items[itemCode])
		} else {
			i++
		}
	}

	character := Character{
		Body:       body,
		Mind:       mind,
		HP:         hp,
		Weapon:     weapon,
		Background: background,
		Items:      itemsGot,
	}

	tmpl := template.Must(template.New("index.html").ParseFiles("index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", character)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
})

func main() {
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.Handle("/", generate).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":"+port, r)
}
