// router/router.go
package router

import (
	"log"
	"net/http"

	"github.com/go-squad-5/quiz-master/internal/app"
	"github.com/go-squad-5/quiz-master/internal/handlers"
)

var currenctNum = 0

type Router struct {
	handlers handlers.Handler
}

// Create router and initialize handlers and concurrency control
func NewRouter(app *app.App) *Router {
	return &Router{
		handlers: handlers.NewHandler(*app.Repository),
	}
}

// Implement ServeHTTP manually
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	// TODO: need to make sure serveHTTP exit only when request is processed
	switch {
	case path == "/quiz/fetch":
		r.handlers.GetQuiz(w, req)
	case path == "/quiz/score":
		log.Println("request recieved to get score and post answers")
		r.handlers.ScoreQuiz(w, req)
	default:
		http.NotFound(w, req)
	}
}

// package router
//
// import (
// 	"net/http"
//
// 	"github.com/go-squad-5/quiz-master/internal/app"
// 	"github.com/go-squad-5/quiz-master/internal/handlers"
// )
//
// type Router struct {
// 	app *app.App
// }
//
// func NewRouter(app *app.App) *Router {
// 	return &Router{
// 		app: app,
// 	}
// }
//
// func (r *Router) Route() *http.ServeMux {
// 	mux := http.NewServeMux()
//
// 	handlers := handlers.NewHandler(*r.app.Repository)
//
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("hello world"))
// 	})
// 	mux.HandleFunc("/quiz/fetch", handlers.GetQuiz)
// 	mux.HandleFunc("/quiz/score", handlers.ScoreQuiz)
//
// 	return mux
// }

// func (r *Router) Route() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	// Initialize the handler with the app instance
// 	h := handlers.NewHandler(r.app)
//
// 	// docs
// 	mux.HandleFunc("/docs/openapi.yaml", h.ServeOpenAPISpec)
// 	mux.Handle("/docs/", h.SwaggerUIHandler())
//
// 	// Versioned routes
// 	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", r.v1Routes()))
//
// 	// add more versioned routes as needed
//
// 	return mux
// }

// func (r *Router) v1Routes() http.Handler {
// 	// Example of adding a versioned route
// 	mux := http.NewServeMux()
//
// 	h := handlers.NewHandler(r.app)
//
// 	// Health check route
// 	mux.HandleFunc("/health", h.HealthCheck)
//
// 	// Add more v1 routes here
//
// 	return mux
// }
