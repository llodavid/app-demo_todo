package main

import (
	"RobertTC32/example-demo_hello/src/commons"
	"RobertTC32/example-demo_hello/src/todo"
	todo_store "RobertTC32/example-demo_hello/src/todo/store"
	todo_ui "RobertTC32/example-demo_hello/src/todo/ui"
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"os"
	"time"
)

// src folder (or special "cmd" subfolder) contains logic to handle the application startup/shutdown,
// and logic to coordinate between business context based modules;
// business context based modules can be used for distinct parts of the business logic with minimal coupling;
// the integration of multiple modules is only needed in more complex applications;
// this simple crud application contains only one module called "todo"
//
// each module has a hexagonal (not layered) structure to minimize technical coupling (and control dependencies):
// - system inside:
// 	Which parts?
// 	"application" (service/facade with executable tasks), "model" (business logic)
// - system outside:
// 	system inside (containing the application logic) needs to be protected from uncontrolled outside (containing external I/O) access;,
// 	using contracts (called 'ports') and technical implementations (called 'adapters');
// 	this makes the system also easier to maintain and easier to test (by using mock implementations of the adapters);
// 	two types of adapters are used:
//  * "inbound"/"driving" adapters (to receive requests):
// 	  eg "store" (data persistence, eg DB, file system), "called api/service" (eg REST API, message queue, email, sms);
// 	* "outbound"/"driven" adapters (to send requests):
//     eg "view/presentation" (e.g. web pages), "calling api/service" (eg REST API, message queue, email, sms);
// this simple application doesn't contain any "api/service" layer;
// it also doesn't define "ports" (interfaces) for the "store" and "ui" adapter, to simplify the code

//go:embed resources
var resourcesEmbed embed.FS

// database with application lifetime
// (needs to be gracefully closed)
var store *todo_store.Store

func main() {
	commons.LoadEnvFile()
	commons.InitLoggerFromEnv()
	slog.Debug("main::main() - Executing")
	//
	// create (and destroy) "store" implementation
	store, err := todo_store.NewStore()
	if err != nil {
		slog.Error("main::main() - Failed to create store", "error", err)
		return
	}
	defer store.Destroy()
	//
	// create (and destroy) "app" implementation with "store" injected
	app, err := todo.NewAppService(store)
	if err != nil {
		slog.Error("main::main() - Failed to create app", "error", err)
		return
	}
	defer app.Destroy()
	//
	// get access to the embedded resources folder
	//resourcesFs := os.DirFS("src/resources")
	resourcesFs, _ := fs.Sub(resourcesEmbed, "resources")
	//
	// create (and destroy) "ui" implementation with "app" injected
	ui, err := todo_ui.NewUiImpl(app, resourcesFs)
	if err != nil {
		slog.Error("main::main() - Failed to create UI", "error", err)
		return
	}
	defer ui.Destroy()
	//
	// start ui using web server
	srv, _ := commons.NewServer(ui.Router)
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	slog.Info("main::main() - Web Server is available at http://" + host + ":" + port)
	slog.Info("main::main() - Press Ctrl+C to stop")
	if err := srv.RunServer(context.Background(), 5*time.Second); err != nil {
		slog.Error("main::main() - Server error", "error", err)
	}
	//
	slog.Info("main::main() - Stopped")
}
