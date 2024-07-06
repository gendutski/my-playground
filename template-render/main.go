package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title        string
	ContentTitle string
	ContentBody  string
	SidebarItems []string
}

func main() {
	// Memparsing template
	tmpl, err := template.ParseFiles("layout.html", "content.html", "sidebar.html")
	if err != nil {
		panic(err)
	}

	// Data untuk template
	data := map[string]interface{}{
		"Title":        "My Website",
		"ContentTitle": "Welcome to My FUckin Website",
		"ContentBody":  "This is the main content of the website.",
		"SidebarItems": []map[string]string{
			{"item": "Item 1"},
			{"item": "Item 2"},
			{"item": "Item 3"},
		},
	}

	// Menjalankan server HTTP
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":1234", nil)
}
