package db

import "errors"


type Player struct {
	ID     string  `json:"playerId"`
	Points float64 `json:"balance"`
}

type dbPlayers struct {
	Data map[string]*Player
}

func (ps *dbPlayers) Get(id string) (*Player, error){
	if p, ok := ps.Data[id]; ok{
		return p, nil
	}
	return nil, errors.New("not found")
}

func (ps *dbPlayers) Set(p *Player) error{
	if p.ID == ""{
		return errors.New("invalid player ID")
	}
	if _, ok := ps.Data[p.ID]; ok{
		return errors.New("duplicate player ID")
	}
	ps.Data[p.ID] = p
	return nil
}

func (ps *dbPlayers) Clear(){
	for k := range ps.Data {
		delete(ps.Data, k)
	}
}

var Players *dbPlayers

func init(){
	Players = &dbPlayers{}
	Players.Data = make(map[string]*Player)
}