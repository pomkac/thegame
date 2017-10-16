package db

import (
	"github.com/hashicorp/go-memdb"
	"log"
	"errors"
)

type Database struct{

}

var (
	E_DB_NOT_FOUND = errors.New("Record not found")
	E_DB_INVALID_ID = errors.New("Invalid ID")
	E_DB_INVALID_VALUE = errors.New("Invalid value")
	M_DB_CREATE_ERROR = "Cannot create database"
)

var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		TablePlayer: &memdb.TableSchema{
			Name: TablePlayer,
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
			},
		},
		"tournament": &memdb.TableSchema{
			Name: "tournament",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
			},
		},
	},
}

var DB *memdb.MemDB

func init(){
	var err error
	DB, err = memdb.NewMemDB(schema)

	if err != nil{
		log.Fatal(M_DB_CREATE_ERROR)
		return
	}
}
