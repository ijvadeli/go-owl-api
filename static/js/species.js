console.log("🦉 | Species JavaScript loaded");

// Fetch API data by creating http request to /species-data
fetch("/species-data")
    // Take response and convert json into js data
    .then(response => {
        // If response is not ok
        if (!response.ok) {
            throw new Error("Failed to fetch owl species");
        }
        // Else just return data
        return response.json();
    })
    // Give js data as species
    .then(species => {
        const speciesList = document.getElementById("species-list");

        // For every owl in species array run:
        species.forEach(owl => {
            // Create card
            const card = document.createElement("div");

            // Change the card info to species data
            card.innerHTML = `
                <h2>${owl.name}</h2>
                <p>Scientific name: ${owl.scientific_name}</p>
                <p>Region: ${owl.region}</p>
                <a href="/species/${owl.id}">View species</a>
            `

            // Append the species data to card
            speciesList.appendChild(card);
        });
    })
    .catch(error => {
        console.error(error);

        document.getElementById("species-list").textContent =
        "Could not load owl species data";
    });