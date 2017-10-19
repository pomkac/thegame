package controllers

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/db"
)

const (
	PlayerTemplate  = `{"playerId": "%s", "balance": %.2f}`
	QueryPlayerId   = "playerId"
	QueryPoints     = "points"
	ContentTypeJSON = "application/json; charset=utf-8"
)

func PlayerFund(ctx *fasthttp.RequestCtx) (err error) {
	query := ctx.QueryArgs()
	respError := fasthttp.StatusOK

	defer func() {
		if respError != fasthttp.StatusOK {
			ctx.SetStatusCode(respError)
			ctx.SetConnectionClose()
		}
	}()

	// Check if "playerId" and "points" exists in query string
	if !query.Has(QueryPlayerId) || !query.Has(QueryPoints) {
		respError = fasthttp.StatusBadRequest
		return
	}

	id := string(query.Peek(QueryPlayerId))
	points, err := query.GetUfloat(QueryPoints)

	// Check if "points" is float and is non-negative number
	if err != nil || points < 0 {
		respError = fasthttp.StatusBadRequest
		return
	}

	player := &db.Player{ID: id}
	conn := db.DB.Conn()
	player.FundTake(points, &conn)
	conn.Close()
	return
}

func PlayerTake(ctx *fasthttp.RequestCtx) (err error) {
	conn := db.DB.Conn()
	query := ctx.QueryArgs()
	respError := fasthttp.StatusOK

	defer func() {
		conn.Close()
		if respError != fasthttp.StatusOK {
			ctx.SetStatusCode(respError)
			ctx.SetConnectionClose()
		}
	}()

	// Check if "playerId" and "points" exists in query string
	if !query.Has(QueryPlayerId) || !query.Has(QueryPoints) {
		respError = fasthttp.StatusBadRequest
		return
	}

	id := string(query.Peek(QueryPlayerId))
	points, err := query.GetUfloat(QueryPoints)

	// Check if "points" is float and is non-negative number
	if err != nil || points < 0 {
		respError = fasthttp.StatusBadRequest
		return
	}

	player, err := db.Players.Get(id, &conn)

	// Check if player exists
	if player.ID == "" {
		respError = fasthttp.StatusNotFound
		return
	}

	err = player.FundTake(-points, &conn)

	if err != nil {
		respError = fasthttp.StatusPaymentRequired
		return
	}
	return
}

func PlayerBalance(ctx *fasthttp.RequestCtx) (err error) {
	query := ctx.QueryArgs()
	respError := fasthttp.StatusOK

	defer func() {
		if respError != fasthttp.StatusOK {
			ctx.SetStatusCode(respError)
			ctx.SetConnectionClose()
		} else {
			ctx.SetContentType(ContentTypeJSON)
		}
	}()

	// Check if "playerId" exists in query string
	if !query.Has(QueryPlayerId) {
		respError = fasthttp.StatusBadRequest
		return
	}

	id := string(query.Peek(QueryPlayerId))

	conn := db.DB.Conn()
	player, err := db.Players.Get(id, &conn)
	conn.Close()

	// Check if player exists
	if err != nil {
		respError = fasthttp.StatusNotFound
		return
	}

	// Format data for "precision 2"
	fmt.Fprintf(ctx, PlayerTemplate, player.ID, player.Points)
	return
}
