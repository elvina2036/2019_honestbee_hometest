package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Page : http page data
type Page struct {
	Title1, Title2, Title3       string
	Content1, Content2, Content3 int
}

const (
	title1 = "Current Connection Count: "
	title2 = "Processed Request Count: "
)

func loadPage() (*Page, error) {
	return &Page{
		Title1:   title1,
		Title2:   title2,
		Content1: connectioncnt,
		Content2: processedreqcnt}, nil
}

func handleHTTPListener() {
	http.HandleFunc("/view/", viewHandler)
	go http.ListenAndServe(":4000", nil)
	fmt.Println("HTTP Server On.")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage()
	t, _ := template.ParseFiles("../view.html")
	t.Execute(w, p)
}
