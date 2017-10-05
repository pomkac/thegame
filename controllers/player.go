package controllers

import (
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/db"
	"fmt"
)

const PLAYER_TEMPLATE = `{"playerId": "%s", "balance": %.2f}`

func PlayerFund(ctx *fasthttp.RequestCtx) {
	// Check if "playerId" and "points" exists in query string
	if !ctx.QueryArgs().Has("playerId") || !ctx.QueryArgs().Has("points"){
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	points, err := ctx.QueryArgs().GetUfloat("points")

	// Check if "points" is float
	if err != nil{
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	// Check if "points" is non-negative number
	if  points < 0{
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	id := string(ctx.QueryArgs().Peek("playerId"))

	player, err := db.Players.Get(id)

	// Check if player exists
	if err != nil{
		// If player not exists then create it
		err = db.Players.Set(&db.Player{id, points})
	}else{
		// Add points to exists player
		player.Points += points
	}
}

func PlayerTake(ctx *fasthttp.RequestCtx) {
	// Check if "playerId" and "points" exists in query string
	if !ctx.QueryArgs().Has("playerId") || !ctx.QueryArgs().Has("points"){
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	id := string(ctx.QueryArgs().Peek("playerId"))

	player, err := db.Players.Get(id)

	// Check if player exists
	if err != nil{
		ctx.SetStatusCode(404)
		ctx.SetConnectionClose()
		return
	}

	points, err := ctx.QueryArgs().GetUfloat("points")

	// Check if "points" is float
	if err != nil{
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	// Check if "points" is non-negative number and current player balance greater or equals points
	if player.Points - points < 0 || points < 0{
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	player.Points -= points
}

func PlayerBalance(ctx *fasthttp.RequestCtx) {

	// Check if "playerId" exists in query string
	if !ctx.QueryArgs().Has("playerId"){
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	id := string(ctx.QueryArgs().Peek("playerId"))

	player, err := db.Players.Get(id)

	// Check if player exists
	if err != nil{
		ctx.SetStatusCode(404)
		ctx.SetConnectionClose()
		return
	}

	ctx.SetContentType("application/json; charset=utf-8")

	// Format data for "precision 2"
	fmt.Fprintf(ctx, PLAYER_TEMPLATE, player.ID, player.Points)
}

