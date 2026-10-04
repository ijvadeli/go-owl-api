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
	About          string `json:"About"`
}

// Create a slice of species using []Species
var speciesData = []Species{
	{
		ID:             1,
		Name:           "Striped Owl",
		ScientificName: "Asio clamator",
		Region:         "South America and Central America",
		Image:          "/static/img/asio-clamator.jpg",
		About:          "Placeholder content",
	},
	{
		ID:             2,
		Name:           "Barn Owl",
		ScientificName: "Tyto alba",
		Region:         "Europe, Africa and West Asia",
		Image:          "/static/img/tyto-alba.png",
		About:          "Placeholder content",
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
