package api

import (
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"net/http"
)

type public struct {
	dataService pubStore
}

type pubStore interface {
	Hello() (helloStr string)
}

func (s *public) sayHello(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, R.JSON{"data": s.dataService.Hello()})
}
