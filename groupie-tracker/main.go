package main

import (
	"log"
	"net/http"
	"text/template"
)

var (
	homeTmpl     *template.Template
	artistTmpl   *template.Template
	error400Tmpl *template.Template
	error404Tmpl *template.Template
	error500Tmpl *template.Template
)

func main() {

	var err error
	homeTmpl, err = template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	artistTmpl, err = template.ParseFiles("static/artist.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	error400Tmpl, err = template.ParseFiles("static/400.html")
	if err != nil {
		log.Fatal("Error parsing 400 template:", err)
	}

	error404Tmpl, err = template.ParseFiles("static/404.html")
	if err != nil {
		log.Fatal("Error parsing 404 template:", err)
	}

	error500Tmpl, err = template.ParseFiles("static/500.html")
	if err != nil {
		log.Fatal("Error parsing 500 template:", err)
	}
	

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist/", ArtistHandler)

	log.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
