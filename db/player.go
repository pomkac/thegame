package db

import (
	"github.com/pomkac/mnemonic"
)

const TablePlayer = "player_"

type Player struct {
	ID     string  `json:"playerId"`
	Points float64 `json:"balance"`
}

type players struct{}

var Players = &players{}

func (p *players) Get(id string, conn *mnemonic.Connection) (Player, error) {
	raw, err := conn.Get(TablePlayer + id)
	if err != nil {
		return Player{}, err
	}

	return raw.(Player), nil
}

func (p *Player) Save(conn *mnemonic.Connection) error {
	if p.ID == "" {
		return errorInvalidId
	}

	conn.Set(TablePlayer+p.ID, *p)

	return nil
}

func (p *Player) FundTake(points float64, conn *mnemonic.Connection) (err error) {
	if p.ID == "" {
		return errorInvalidId
	}

	raw, err := conn.Get(TablePlayer + p.ID)
	if err == nil {
		*p = raw.(Player)
	}

	if p.Points+points < 0 {
		return errorInvalidValue
	}

	p.Points += points

	p.Save(conn)

	return nil
}
