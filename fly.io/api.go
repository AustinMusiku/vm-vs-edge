package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var characters = []string{
	"Rick Grimes",
	"Morgan Jones",
	"Duane Jones",
	"Jenny Jones",
	"Carl Grimes",
	"Lori Grimes",
	"Shane Walsh",
	"Daryl Dixon",
	"Merle Dixon",
	"Dale Horvath",
	"Andrea Harrison",
	"Ed Peletier",
	"Carol Peletier",
	"Sophia Peletier",
	"Glenn Rhee",
	"T-Dog",
	"Jacqui",
	"Jim",
	"Amy",
	"Luis Morales",
	"Miranda Morales",
	"Eliza Morales",
	"Edward Jenner",
	"Maggie Greene",
	"Hershel Greene",
	"Beth Greene",
	"Patricia",
	"Jimmy",
	"Otis",
	"Randall Culver",
	"Judith Grimes",
	"Michonne",
	"The Governor",
	"Tyreese Williams",
	"Sasha Williams",
	"Karen",
	"Ceasar Martinez",
	"Shumpert",
	"Oscar (Big Tiny)",
	"Axel",
	"Big Tiny",
	"Tomas",
	"Andrew",
	"Haley",
	"Milton Mamet",
	"Dr. Stevens",
	"Tara Chambler",
	"Ryan Samuels",
	"Lizzie Samuels",
	"Mika Samuels",
	"Abraham Ford",
	"Eugene Porter",
	"Rosita Espinosa",
	"Bob Stookey",
	"Gareth",
	"Mary",
	"Chris",
	"Martin",
	"David",
	"Greg",
	"Mike",
	"Albert",
	"Charlie",
	"Pete Dolgen",
	"Gabriel Stokes",
	"Noah",
	"Dawn Lerner",
	"Steven Edwards",
	"O'Donnell",
	"Michael Todd",
	"Aaron",
	"Eric Raleigh",
	"Jessie Anderson",
	"Ron Anderson",
	"Sam Anderson",
	"Deanna Monroe",
	"Spencer Monroe",
	"Aiden Monroe",
	"Reg Monroe",
	"Nicholas",
	"Porcine",
	"Shelly Neudermeyer",
	"Sturgess",
	"Scott",
	"Heath",
	"Olivia",
	"Denise Cloyd",
	"Nicholas",
	"Enid",
	"Gregory",
	"Paul Rovia",
	"Simon",
	"Negan Smith",
	"Dwight",
	"Jerry",
	"Ezekiel",
	"Shiva",
	"Benjamin",
	"Richard",
	"Jadis",
	"Tamiel",
	"Brion",
	"Kurt",
	"Sherry",
	"Arat",
	"Cindy",
	"Beatrice",
	"Kathy",
	"Gary",
	"Paula",
	"Molly",
	"Harlan Carson",
	"Alden",
	"Alpha",
	"Beta",
	"Lydia",
	"Laura",
	"Luke",
	"Connie",
	"Kelly",
	"Magna",
	"Yumiko",
	"Lucille Smith",
	"Virgil",
	"Leah Shaw",
	"Sebastian Milton",
	"Agatha",
	"Max",
	"Roy",
	"Tommy",
	"Marcus",
	"Rufus",
	"Devin",
	"Wayne",
	"Jace",
	"Kevin",
	"Ruthie",
	"Kimber",
	"Joseph",
	"Jen",
	"Tyler",
	"Adeline",
	"Emerson",
	"Princess",
	"Mercer",
	"Sebastian",
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleRenderIndex)
	mux.HandleFunc("GET /characters", handleListCharacters)

	http.ListenAndServe(":8080", mux)
}

func handleRenderIndex(w http.ResponseWriter, r *http.Request) {
	listItems := ""
	for _, s := range characters {
		listItems += fmt.Sprintf("<li>%s</li>", s)
	}

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<meta name="description" content="The Walking Dead characters">

			<title>The Walking Dead characters</title>
		</head>
		<body>
			<h1>The Walking Dead characters</h1>
			<p>Here are the characters from The Walking Dead:</p>
			<ol>%s</ol>
		</body>
	</html>
	`, listItems)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleListCharacters(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
