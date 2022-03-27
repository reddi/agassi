package service

import (
	"github.com/reddi/agassi/app/store"
	"github.com/reddi/agassi/app/store/engine"
)

type DataStore struct {
	Engine engine.Interface
}

func (s *DataStore) CreatePlayer(player store.Player) (playerID string, err error) {
	return s.Engine.CreatePlayer(player)
}

func (s *DataStore) CreateCoach(coach store.Coach) (coachID string, err error) {
	return s.Engine.CreateCoach(coach)
}

func (s *DataStore) Close() error {
	return s.Engine.Close()
}
