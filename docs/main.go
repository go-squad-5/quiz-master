// package main
//
// import (
//
//	"log"
//	"net/http"
//
//	// "github.com/go-squad-5/quiz-master/internal/app"
//	// "github.com/go-squad-5/quiz-master/internal/config"
//	"github.com/go-squad-5/quiz-master/internal/router"
//
// )
//
//	func main() {
//		// Load configuration
//		// cfg := config.Load()
//		// Initialize the application
//		// app := app.NewApp(cfg)
//		// Initialize the router with the app instance
//		route := router.Route()
//		// Start the HTTP server
//		log.Println("Starting server on :8080")
//		log.Fatal(http.ListenAndServe(":8080", route))
//	}
package main

import (
	"log"
	"net/http"
	"path/filepath"
)

type Handler struct{}

func (h *Handler) ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("docs", "openapi.yaml"))
}

func (h *Handler) SwaggerUIHandler() http.Handler {
	return http.StripPrefix("/docs/", http.FileServer(http.Dir("docs/swagger")))
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Route() *http.ServeMux {
	mux := http.NewServeMux()
	// Initialize the handler with the app instance
	h := &Handler{}

	// docs
	mux.HandleFunc("/docs/openapi.yaml", h.ServeOpenAPISpec)
	mux.Handle("/docs/", h.SwaggerUIHandler())

	return mux
}

func main() {
	route := NewRouter().Route()
	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Println("view docs on http://localhost:8080/docs/")
	log.Fatal(http.ListenAndServe(":8080", route))
}
