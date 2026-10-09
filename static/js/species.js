console.log("🦉 | Species JavaScript loaded");

// Create search items
const speciesList = document.getElementById("species-list");
const searchInput = document.getElementById("species-search");
const noResults = document.getElementById("no-results");

// Store species after fetching the API
let allSpecies = [];

// Fetch API data by creating http request to /species-data
fetch("/species-data")
  // Take response and convert json into js data
  .then((response) => {
    // If response is not ok
    if (!response.ok) {
      throw new Error("Failed to fetch owl species");
    }
    // Else just return data
    return response.json();
  })
  // Give js data as species
  .then((species) => {
    allSpecies = species;

    displaySpecies(allSpecies);

    // For every owl in species array run:
  })
  .catch((error) => {
    console.error(error);

    speciesList.textContent = "Could not load owl species data";
  });

// Display list of owls
function displaySpecies(species) {
  speciesList.replaceChildren();

  noResults.hidden = species.length !== 0;

  species.forEach((owl) => {
    // Create card
    const card = document.createElement("a");
    card.className = "species-card";
    card.href = `/species/${owl.id}`;

    const heading = document.createElement("h2");
    heading.textContent = `${owl.name} (${owl.scientific_name})`;

    const image = document.createElement("img");
    image.className = "owl-img";
    image.src = owl.image;
    image.alt = owl.name;

    // Change the card info to species data
    card.append(
        heading,
        image
    );

    // Append the species data to card
    speciesList.appendChild(card);
  });
}

// Filter list whenever user types
searchInput.addEventListener("input", () => {
    const searchTerm = searchInput.value.trim().toLowerCase();

    const filteredSpecies = allSpecies.filter(owl => {
        const searchableText = [
            owl.name,
            owl.scientific_name,
            owl.region,
            owl.description
        ].join(" ").toLowerCase();

        return searchableText.includes(searchTerm);
    })
    
    displaySpecies(filteredSpecies);
})
