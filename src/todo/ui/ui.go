package todo_ui

import (
	todo "RobertTC32/example-demo_hello/src/todo"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
)

// "ui" (aka views and controllers) contains adapter for web gui in the "Hexagonal Architecture"

type Ui struct {
	app         *todo.AppService
	resourcesFs fs.FS
	Router      *http.ServeMux
}

// prevent parsing template files on every request by caching them;
// only used in listHandler example where template library is used
var templates = make(map[string]*template.Template)

func NewUiImpl(app *todo.AppService, resourcesFs fs.FS) (*Ui, error) {
	slog.Debug("ui::NewUiImpl() - Executing")
	router := http.NewServeMux()
	ui := &Ui{
		app:         app,
		resourcesFs: resourcesFs,
		Router:      router,
	}
	//
	// add todo handlers
	router.Handle("GET /public/", http.FileServerFS(resourcesFs))
	router.HandleFunc("GET /", ui.defaultHandler)
	router.HandleFunc("GET /todo", ui.todoHandler)
	router.HandleFunc("POST /todo", ui.todoPostHandler)
	//
	// add extra example handlers
	router.HandleFunc("GET /list", ui.listHandler)
	router.HandleFunc("GET /mapping", ui.mappingHandler)
	router.HandleFunc("GET /hello", ui.helloHandler)
	// prevent parsing template files on every request by caching them
	if len(templates) == 0 {
		templates["list.html"] = template.Must(template.ParseFS(resourcesFs, "templates/list.html"))
	}
	//
	return ui, nil
}

func (t *Ui) Destroy() error {
	return nil
}

func (t *Ui) defaultHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ui::defaultHandler() - Executing")
	http.Redirect(w, r, "/todo", http.StatusMovedPermanently)
}
