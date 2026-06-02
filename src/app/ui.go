package app

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Ui struct {
	router *http.ServeMux
	store  *Store
}

func NewUi(router *http.ServeMux, s *Store) (*Ui, error) {
	slog.Debug("ui::NewUi() - Executing")
	ui := &Ui{router: router, store: s}
	router.HandleFunc("/", ui.getTodosHandler)
	return ui, nil
}

func (ui *Ui) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::getTodosHandler() - Executing")
	slog.Debug("ui::getTodosHandler() - Reading data", "storage", ui.store.dsn)
	d, err := ui.store.GetTodos()
	if err != nil {
		slog.Error("ui::getTodosHandler() - Failed to access data", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//
	//ui.store.GetMappingtests1()
	//ui.store.GetMappingtests2()
	//
	slog.Debug("ui::getTodosHandler() - Building response")
	tmpl := template.Must(template.ParseFiles("./resources/list.html"))
	err = tmpl.Execute(w, d)
	if err != nil {
		slog.Error("ui::getTodosHandler() - Failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
