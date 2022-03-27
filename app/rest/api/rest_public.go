package api

import (
	"github.com/go-chi/render"
	"net/http"
)

type public struct {
	dataService pubStore
}

type pubStore interface {
	Hello() (helloStr string)
}

func (s *public) sayHello(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, s.dataService.Hello())
}
