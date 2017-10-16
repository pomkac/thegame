package controllers

import (
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/db"
	"fmt"
)

const (
	PlayerTemplate  = `{"playerId": "%s", "balance": %.2f}`
	QueryPlayerId   = "playerId"
	QueryPoints     = "points"
	ContentTypeJSON = "application/json; charset=utf-8"
)

func PlayerFund(ctx *fasthttp.RequestCtx) {
	// Check if "playerId" and "points" exists in query string
	if !ctx.QueryArgs().Has(QueryPlayerId) || !ctx.QueryArgs().Has(QueryPoints) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetConnectionClose()
		return
	}

	points, err := ctx.QueryArgs().GetUfloat(QueryPoints)

	// Check if "points" is float
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetConnectionClose()
		return
	}

	// Check if "points" is non-negative number
	if points < 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetConnectionClose()
		return
	}

	id := string(ctx.QueryArgs().Peek(QueryPlayerId))

	player := &db.Player{ID: id}

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetConnectionClose()
		return
	}

	player.FundTake(points)
}

func PlayerTake(ctx *fasthttp.RequestCtx) {
	// Check if "playerId" and "points" exists in query string

	respError := fasthttp.StatusOK

	defer func() {
		if respError != fasthttp.StatusOK {
			ctx.SetStatusCode(respError)
			ctx.SetConnectionClose()
		}
	}()

	if !ctx.QueryArgs().Has(QueryPlayerId) || !ctx.QueryArgs().Has(QueryPoints) {
		respError = fasthttp.StatusBadRequest
		return
	}

	id := string(ctx.QueryArgs().Peek(QueryPlayerId))

	points, err := ctx.QueryArgs().GetUfloat(QueryPoints)

	// Check if "points" is float and is non-negative number
	if err != nil || points < 0 {
		respError = fasthttp.StatusBadRequest
		return
	}

	player, err := db.Players.Get(id)

	// Check if player exists
	if player == nil {
		respError = fasthttp.StatusNotFound
		return
	}

	err = player.FundTake(-points)

	if err != nil {
		respError = fasthttp.StatusPaymentRequired
		return
	}
}

func PlayerBalance(ctx *fasthttp.RequestCtx) {

	// Check if "playerId" exists in query string
	if !ctx.QueryArgs().Has(QueryPlayerId) {
		ctx.SetStatusCode(ERROR_BAD_REQUEST)
		ctx.SetConnectionClose()
		return
	}

	id := string(ctx.QueryArgs().Peek(QueryPlayerId))

	player, err := db.Players.Get(id)

	// Check if player exists
	if err != nil {
		ctx.SetStatusCode(ERROR_NOT_FOUND)
		ctx.SetConnectionClose()
		return
	}

	ctx.SetContentType(ContentTypeJSON)

	// Format data for "precision 2"
	fmt.Fprintf(ctx, PlayerTemplate, player.ID, player.Points)
}
