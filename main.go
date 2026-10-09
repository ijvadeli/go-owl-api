package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// Create owl species struct
type Species struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ScientificName string `json:"scientific_name"`
	Region         string `json:"region"`
	Image          string `json:"image"`
	About          string `json:"about"`
	Description    string `json:"description"`
	History        string `json:"history"`
}

// Create a slice of species using []Species
var speciesData = []Species{
	{
		ID:             1,
		Name:           "Striped owl",
		ScientificName: "Asio clamator",
		Region:         "South and Central America",
		Image:          "/static/img/asio-clamator.jpg",
		About:          "The striped owl, also known as the Asio clamator is a medium 	sized owl species native to South and Central America.",
		Description:    "The striped owl is a relatively large species with recognisable tufts of feathers on his head, resembling ears. It's around 30-38 cm (12-15 inch) in length and weighs from 320-546 g (11.3-19.3 oz). It's head, back, wings and tail all have brown with black stripes and small markings, the underparts are more buff-colored with heavy black streaking on the breast. The face is white with a thin black border around it.",
		History:        "The striped owl was originally described by the French ornithologist Louis Pierre Vieillot in the year 1808, he named it the Bubo clamator. The name clamator is the Latin meaning for shouter. The type locality is Cayenne in French Guiana. The striped owl was at one time placed in it's own genus Rhinoptynx, and was later transferred to the genus Pseudoscops. A molecular study that compared the mitochondrial DNA sequences indicated that it should be placed in the genus Asio instead. This was confirmed by a large study of the owls in 2019.",
	},
	{
		ID:             2,
		Name:           "Barn owl",
		ScientificName: "Tyto alba",
		Region:         "Europe, Africa and West Asia",
		Image:          "/static/img/tyto-alba.png",
		About:          "The barn owls (not to be confused with the Barred owl) are owls in the genus of Tyto, the most common genus of owls in the world.",
		Description:    "The barn owl is a medium-sized owl with a large head and recognizable heart-shaped face. They have long yet strong legs with powerful talons. The term barn owl is often used to describe the whole family Tytonidae, which also includes bay owls in the genus Phodilus. Tyto being the largets genus of Tytonidae including species like the Western, American and Eastern barn owl. They are often referred to as a single species known as barn owl or common barn owl. Then the others are Andaman masked owl, only really seen on the southern Andaman Islands. The New Caledonian barn owl, which went extinct but was seen on the islands of New Caledonia in Melanesia. The Rivero's barn owl, in Cuba which also went extinct.",
		History:        "The barn owl family (Tytonidae) originated in Asia over 45 million years ago, before spreading globally into distinct species. The western barn owl was formally described and names Strix alba in 1769 by the Italian naturalist Giovanni Antonio Scopoli, later moved to the genus Tyto in 1828. Around 3,500-5,500 years ago barn owls began moving from natural cavities into human buildings. For example barns, haystacks and church steeples. They were often associated with superstition and fear, During the Renaissance the perception shifted however. Artists like Michelangelo sculpted barn owls turning the bird into more of a symbol of wisdom and beauty. During the 19th century in Europe, the barn owl populations greatly suffered. Gamekeepers trappen and shot them over false fears of lost gamebirds, and Victorian fashion trades prized their feathers. Population of the barn owls sharply declined through the mid-20th century due to modern building, loss of hollow trees, intensive farming and pesticides. Today, safe artificial nesting boxes managed by groups like the Barn Owl Trust help stabilize numbers.",
	},
}

// * Main function
func main() {
	// Serve static files from static directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// Route handler for home page
	http.HandleFunc("/", home)
	// Route handler for species page
	http.HandleFunc("/species/", species)
	// Route handler for API
	http.HandleFunc("/species-data", speciesAPI)
	http.HandleFunc("/species-data/", speciesAPI)

	// Simple terminal log to show server is running on port 8080
	fmt.Println("🦉 | Server is running on http://localhost:8080")

	// Start server on port 8080 and listen to incoming requests
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Global scope
// Home function handles request to index.html
func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

// Handles request to species page
func species(w http.ResponseWriter, r *http.Request) {
	// Remove trailing / from URL
	path := strings.TrimSuffix(r.URL.Path, "/")
	// Split path into parts
	parts := strings.Split(path, "/")

	// Display on /species
	if len(parts) == 2 {
		renderTemplate(w, "species.html", nil)
		return
	}

	// Display on /species/number
	if len(parts) == 3 {
		id, err := strconv.Atoi(parts[2])
		// If id is invalid (err is not nil)
		if err != nil {
			http.Error(w, "Invalid species ID", http.StatusBadRequest)
			return
		}

		// For all owls in speciesData render data on species-detail/id
		for _, owl := range speciesData {
			if owl.ID == id {
				renderTemplate(w, "species-detail.html", owl)
				return
			}
		}

		// If species is not found
		http.Error(w, "Species not found", http.StatusNotFound)
		return
	}

	// General error message
	http.Error(w, "Not found", http.StatusNotFound)
}

// * SpeciesAPI function
// Handles GET /species-data
func speciesAPI(w http.ResponseWriter, r *http.Request) {
	// If method is not GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set display/content type to json data
	w.Header().Set("Content-Type", "application/json")

	// Remove trailing / from URL
	path := strings.TrimSuffix(r.URL.Path, "/")

	// Split path into parts
	parts := strings.Split(path, "/")

	// GET /species-data
	if len(parts) == 2 {
		json.NewEncoder(w).Encode(speciesData)
		return
	}

	// GET /species-data/1
	if len(parts) == 3 {
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid species ID", http.StatusBadRequest)
			return
		}

		for _, owl := range speciesData {
			if owl.ID == id {
				json.NewEncoder(w).Encode(owl)
				return
			}
		}

		http.Error(w, "Species not found", http.StatusNotFound)
		return
	}

	// r.URL.Path contains path client requested
	// if r.URL.Path == "/species-data" {
	// 	json.NewEncoder(w).Encode(speciesData)
	// 	return
	// }

	// json.NewEncoder(w).Encode(speciesData)
}

// * Render HTML template
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parsing specified template file being passed as input
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Give data and give it to template
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
