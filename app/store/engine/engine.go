package engine

import "github.com/reddi/agassi/app/store"

type Interface interface {
	CreatePlayer(player store.Player) (playerID string, err error)
	CreateCoach(coach store.Coach) (coachID string, err error)
	ListPlayers() (players []store.Player, err error)
	ListCoaches() (coaches []store.Coach, err error)

	AddReview(playerID string, review store.Review) error
	ListReviews(playerID string) ([]store.Review, error)

	Close() error // close storage engine
}
