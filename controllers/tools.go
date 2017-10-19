package controllers

import (
	"errors"
	"strings"
	"encoding/json"
	"github.com/pomkac/thegame/db"
	"github.com/valyala/fasthttp"
)

const (
	strContentType = "Content-Type"
	strJSON        = "application/json"
)

var errorInvalidData = errors.New("invalid data")

func Reset(_ *fasthttp.RequestCtx) error {
	conn := db.DB.Conn()
	conn.Drop()
	conn.Close()
	return nil
}

func ValidateResult(ctx *fasthttp.RequestCtx) (res *ResultTournamentStruct, err error) {
	ct := ctx.Request.Header.Peek(strContentType)

	// Check content type
	if !strings.Contains(string(ct), strJSON) {
		err = errorInvalidData
		return
	}

	data := ctx.PostBody()

	res = &ResultTournamentStruct{}

	// Try parse POST body to struct
	err = json.Unmarshal(data, res)

	if err != nil {
		return
	}

	// Validate struct
	if res.Winners == nil || res.ID == nil {
		err = errorInvalidData
		return
	}

	for _, w := range *res.Winners {
		if w.ID == nil || w.Prize == nil || *w.Prize < 0 {
			err = errorInvalidData
			return
		}
	}

	return
}
