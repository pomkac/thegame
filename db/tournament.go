package db

import (
	"github.com/pomkac/mnemonic"
)

const TableTournament = "tournament_"

type Tournament struct {
	ID      string
	Deposit float64
	Players map[string]*TournamentPlayer
	Ended   bool
}

type TournamentPlayer struct {
	ID     string
	Bakers []Player
}

type tournaments struct{}

var Tournaments = &tournaments{}

func (t *tournaments) Get(id string, conn *mnemonic.Connection) (Tournament, error) {
	raw, err := conn.Get(TableTournament + id)
	if err != nil {
		return Tournament{}, err
	}

	return raw.(Tournament), nil
}

func (t *Tournament) Save(conn *mnemonic.Connection) error {
	if t.ID == "" {
		return errorInvalidId
	}

	conn.Set(TableTournament+t.ID, *t)

	return nil
}
