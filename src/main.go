package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandleFunc)
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Web Server is available at http://localhost:8080/")
	fmt.Println("Press Ctrl+C to stop")
	srv.ListenAndServe()
}

type HelloPageData struct {
	FullName  string
	CurrentDT string
}

func helloHandleFunc(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./src/templates/hello.html"))
	d := HelloPageData{FullName: "Everybody", CurrentDT: time.Now().Format(time.RFC3339)}
	err := tmpl.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
