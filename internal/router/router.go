// router/router.go
package router

import (
	"net/http"

	"github.com/go-squad-5/quiz-master/internal/handlers"
)

var currenctNum = 0

type Router struct {
	handlers handlers.Handler
}

func NewRouter(handler handlers.Handler) *Router {
	return &Router{
		handlers: handler,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	switch path {
	case "/quiz/fetch":
		r.handlers.GetQuiz(w, req)
	case "/quiz/score":
		r.handlers.ScoreQuiz(w, req)
	default:
		http.NotFound(w, req)
	}
}
