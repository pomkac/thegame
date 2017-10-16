package db

import "errors"


const TablePlayer = "player"

type Player struct {
	ID     string  `json:"playerId"`
	Points float64 `json:"balance"`
}

type players struct{}


var Players = &players{}


func (p *players) Get(id string)  (*Player, error){
	txn := DB.Txn(false)
	defer txn.Abort()
	raw, err := txn.First(TablePlayer, "id", id)
	if err != nil {
		return &Player{}, err
	}

	if raw == nil{
		return nil, E_DB_NOT_FOUND
	}
	return raw.(*Player), nil
}

func (p *players) DeleteAll() error{
	txn := DB.Txn(true)
	if _, err := txn.DeleteAll(TablePlayer, "id"); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (p *Player) Save() error{
	if p.ID == ""{
		return errors.New("invalid player ID")
	}
	txn := DB.Txn(true)
	if err := txn.Insert(TablePlayer, p); err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}

func (p *Player) FundTake(points float64) (err error){
	if p.ID == ""{
		return E_DB_INVALID_ID
	}
	txn := DB.Txn(true)

	defer func(){
		if err != nil{
			txn.Abort()
			return
		}
		txn.Commit()
	}()

	raw, err := txn.First(TablePlayer, "id", p.ID)
	if err != nil {
		return err
	}

	if raw != nil{
		*p = *raw.(*Player)
	}

	if p.Points + points < 0{
		return 	E_DB_INVALID_VALUE
	}

	p.Points += points

	if err := txn.Insert(TablePlayer, p); err != nil {
		return err
	}

	return nil
}


