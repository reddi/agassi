package api

import (
	"github.com/go-chi/render"
	log "github.com/go-pkgz/lgr"
	"github.com/reddi/agassi/app/rest"
	"github.com/reddi/agassi/app/store"
	"net/http"
)

type public struct {
	dataService pubStore
}

type pubStore interface {
	CreatePlayer(player store.Player) (playerID string, err error)
	CreateCoach(player store.Coach) (coachID string, err error)
}

func (s *public) createPlayerCtrl(w http.ResponseWriter, r *http.Request) {
	player := store.Player{}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &player); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind player", rest.ErrDecode)
		return
	}
	id, err := s.dataService.CreatePlayer(player)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't save player", rest.ErrInternal)
		return
	}
	player.ID = id
	log.Printf("[DEBUG] created player %+v", player)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &player)
}

func (s *public) createCoachCtrl(w http.ResponseWriter, r *http.Request) {
	coach := store.Coach{}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &coach); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind coach", rest.ErrDecode)
		return
	}
	id, err := s.dataService.CreateCoach(coach)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't save coach", rest.ErrInternal)
		return
	}
	coach.ID = id
	log.Printf("[DEBUG] created coach %+v", coach)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &coach)
}
