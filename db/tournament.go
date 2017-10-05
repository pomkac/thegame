package db

import "errors"

type Tournament struct {
	ID     string
	Deposit float64
	Players map[string]*TournamentPlayer
	Ended bool
}

type TournamentPlayer struct {
	ID string
	Bakers []*Player
}

type dbTournaments struct {
	Data map[string]*Tournament
}

func (ps *dbTournaments) Get(id string) (*Tournament, error){
	if t, ok := ps.Data[id]; ok{
		return t, nil
	}
	return nil, errors.New("not found")
}

func (ps *dbTournaments) Set(t *Tournament) error{
	if t.ID == ""{
		return errors.New("invalid ID")
	}
	if _, ok := ps.Data[t.ID]; ok{
		return errors.New("duplicate ID")
	}
	ps.Data[t.ID] = t
	return nil
}

func (ps *dbTournaments) Clear(){
	for k := range ps.Data {
		delete(ps.Data, k)
	}
}

var Tournaments *dbTournaments

func init(){
	Tournaments = &dbTournaments{}
	Tournaments.Data = make(map[string]*Tournament)
}
