package controllers

import (
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/db"
	"runtime/debug"
	"strings"
	"encoding/json"
	"errors"
)

const (
	INVALID_DATA           = "invalid data"
	ERROR_BAD_REQUEST      = 400
	ERROR_PAYMENT_REQUIRED = 402
	ERROR_FORBIDDEN        = 403
	ERROR_NOT_FOUND        = 404
	ERROR_NOT_ACCEPTABLE   = 406
)

func Reset(ctx *fasthttp.RequestCtx) {
	// Clear Players DB
	db.Players.Clear()
	// Clear Tournaments DB
	db.Tournaments.Clear()
	// Clear memory
	debug.FreeOSMemory()
}

func ValidateResult(ctx *fasthttp.RequestCtx) (res *ResultTournamentStruct, err error) {
	ct := ctx.Request.Header.Peek("Content-Type")

	// Check content type
	if !strings.Contains(string(ct), "application/json") {
		err = errors.New(INVALID_DATA)
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
		err = errors.New(INVALID_DATA)
		return
	}

	for _, w := range *res.Winners {
		if w.ID == nil || w.Prize == nil || *w.Prize < 0 {
			err = errors.New(INVALID_DATA)
			return
		}
	}

	return
}
