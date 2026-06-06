package todo_ui

import (
	"log/slog"
	"net/http"
	"strconv"
)

// "ui" (aka views and controllers) contains adapter for web gui in the "Hexagonal Architecture"

func (t *Ui) todoHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ctrl::todoHandler() - Executing")
	// processing input
	//
	// domain logic
	todos, err := t.app.GetTodos()
	if err != nil {
		slog.Error("ctrl::todoHandler() - Failed to execute logic", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// processing output (using templ)
	err = TodoPage(todos).Render(r.Context(), w)
	if err != nil {
		slog.Error("ctrl::todoHandler() - Failed to process output", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Ui) todoPostHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ctrl::todoPostHandler() - Executing")
	// processing input
	err := r.ParseForm()
	if err != nil {
		slog.Error("ctrl::todoPostHandler() - Failed to process input", "error", err)
		// error can be improved for view
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.Form.Get("id")
	version := r.Form.Get("version")
	action_delete := r.Form.Get("action_delete")
	action_update := r.Form.Get("action_update")
	title := r.Form.Get("title")
	action_insert := r.Form.Get("action_insert")
	//header := r.Header.Values("Content-Type")
	//slog.Debug("ctrl::todoPostHandler() - Form values", "header", header, "id", id, "version", version, "title", title)
	//slog.Debug("ctrl::todoPostHandler() - Form actions", "insert", action_insert, "update", action_update, "delete", action_delete)
	//
	// domain logic
	if len(action_insert) > 0 {
		slog.Debug("ctrl::todoPostHandler() - Inserting", "title", title)
		err = t.app.CreateTodoTitle(title)
		if err != nil {
			slog.Error("ctrl::todoPostHandler() - Failed to execute logic", "error", err)
			// error can be improved for view
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(action_update) > 0 {
		idNum, _ := strconv.ParseUint(id, 10, 64)
		versionNum, _ := strconv.ParseUint(version, 10, 64)
		slog.Debug("ctrl::todoPostHandler() - Updating", "id", idNum, "version", versionNum)
		err = t.app.UpdateTodoCompleted(uint32(idNum), uint32(versionNum))
		if err != nil {
			slog.Error("ctrl::todoPostHandler() - Failed to execute logic", "error", err)
			// error can be improved for view
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(action_delete) > 0 {
		idNum, _ := strconv.ParseUint(id, 10, 64)
		versionNum, _ := strconv.ParseUint(version, 10, 64)
		slog.Debug("ctrl::todoPostHandler() - Deleting", "id", idNum, "version", versionNum)
		err = t.app.DeleteTodo(uint32(idNum), uint32(versionNum))
		if err != nil {
			slog.Error("ctrl::todoPostHandler() - Failed to execute logic", "error", err)
			// error can be improved for view
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	todos, err := t.app.GetTodos()
	if err != nil {
		slog.Error("ctrl::todoPostHandler() - Failed to execute logic", "error", err)
		// error can be improved for view
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// processing output (using templ)
	err = TodoList(todos).Render(r.Context(), w)
	if err != nil {
		slog.Error("ctrl::todoPostHandler() - Failed to process output", "error", err)
		// error can be improved for view
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
