package todo_ui

import (
	"fmt"
	"log/slog"
	"net/http"
)

// "ui" (aka views and controllers) contains adapter for web gui in the "Hexagonal Architecture"

// easy template example
func (t *Ui) listHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ctrl::listHandler() - Executing")
	// processing input
	//
	// domain logic
	slog.Debug("ctrl::listHandler() - Reading data")
	d, err := t.app.GetTodos()
	if err != nil {
		slog.Error("ctrl::listHandler() - Failed to execute logic", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// processing output (using template)
	slog.Debug("ctrl::listHandler() - Building response")
	tmpl := templates["list.html"]
	err = tmpl.Execute(w, d)
	if err != nil {
		slog.Error("ctrl::listHandler() - Failed to process output", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// example to show mapping between mysql/mariadb and go data types (including "null")
func (t *Ui) mappingHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ctrl::dbmappingHandler() - Executing")
	// processing input
	//
	// domain logic
	t.app.GetDbmapping1()
	t.app.GetDbmapping2()
	//
	// processing output (using fmt)
	fmt.Fprintf(w, "Results of dbmapping shown in console")
}

// easy templ example
func (t *Ui) helloHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ctrl::helloTemplHandler() - Executing")
	// processing input
	//
	// domain logic
	//
	// processing output (using templ)
	err := HelloPage("Robert The Coder").Render(r.Context(), w)
	if err != nil {
		slog.Error("ctrl::helloHandler() - Failed to process output", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
