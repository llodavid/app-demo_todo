package app

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Ui struct {
	router  *http.ServeMux
	storage *Storage
}

func NewUi(router *http.ServeMux, s *Storage) (Ui, error) {
	ui := Ui{router: router, storage: s}
	router.HandleFunc("/", ui.listHandleFunc)
	return ui, nil
}

func (ui Ui) listHandleFunc(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::listHandleFunc() - Started")
	slog.Debug("ui::listHandleFunc() - Reading data", "storage", ui.storage.dsn)
	d, err := ui.storage.FindAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//
	slog.Debug("ui::listHandleFunc() - Building response")
	tmpl := template.Must(template.ParseFiles("./resources/list.html"))
	err = tmpl.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
