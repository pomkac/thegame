package controllers

import (
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/db"
)

type ResultTournamentStruct struct {
	ID      *string   `json:"tournamentId"`
	Winners *[]Winner `json:"winners"`
}

type Winner struct {
	ID    *string  `json:"playerId"`
	Prize *float64 `json:"prize"`
}

func AnnounceTournament(ctx *fasthttp.RequestCtx) {
	// Check if "tournamentId" and "deposit" exists in query string
	if !ctx.QueryArgs().Has("tournamentId") || !ctx.QueryArgs().Has("deposit") {
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	deposit, err := ctx.QueryArgs().GetUfloat("deposit")

	// Check if "deposit" is float
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	// Check if "deposit" is non-negative number
	if deposit < 0 {
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	id := string(ctx.QueryArgs().Peek("tournamentId"))

	// Create tournament
	tournament := &db.Tournament{ID: id, Deposit: deposit}
	tournament.Players = make(map[string]*db.TournamentPlayer)

	// Try save tournament
	err = db.Tournaments.Set(tournament)

	// Return error if tournament exists
	if err != nil {
		ctx.SetStatusCode(406)
		ctx.SetConnectionClose()
		return
	}
}

func JoinTournament(ctx *fasthttp.RequestCtx) {
	// Check if "tournamentId" and "playerId" exists in query string
	if !ctx.QueryArgs().Has("tournamentId") || !ctx.QueryArgs().Has("playerId") {
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	tournamentId := string(ctx.QueryArgs().Peek("tournamentId"))
	playerId := string(ctx.QueryArgs().Peek("playerId"))
	backerIds := ctx.QueryArgs().PeekMulti("backerId")

	tournament, err := db.Tournaments.Get(tournamentId)

	// Return error if tournament not exists
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.SetConnectionClose()
		return
	}

	// Return error if tournament not ended
	if tournament.Ended {
		ctx.SetStatusCode(403)
		ctx.SetConnectionClose()
		return
	}

	// Return error if player exists in current tournament
	if _, ok := tournament.Players[playerId]; ok{
		ctx.SetStatusCode(406)
		ctx.SetConnectionClose()
		return
	}

	player, err := db.Players.Get(playerId)

	// Return error if player not exists
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.SetConnectionClose()
		return
	}

	// Create tournament player
	tPlayer := &db.TournamentPlayer{ID: player.ID}

	tPlayer.Bakers = make([]*db.Player, len(backerIds)+1)
	tPlayer.Bakers[0] = player

	// Find and add backers to player
	if len(backerIds) > 0 {
		for i := 0; i < len(backerIds); i++ {
			backer, err := db.Players.Get(string(backerIds[i]))
			// Return error if backer not exists
			if err != nil {
				ctx.SetStatusCode(404)
				ctx.SetConnectionClose()
				return
			}

			tPlayer.Bakers[i+1] = backer
		}
	}

	// Calculate bid
	bid := tournament.Deposit / float64(len(tPlayer.Bakers))

	// Check player balance
	for i := 0; i < len(tPlayer.Bakers); i++ {
		if tPlayer.Bakers[i].Points < bid {
			ctx.SetStatusCode(402)
			ctx.SetConnectionClose()
			return
		}
	}

	// Create bid
	for i := 0; i < len(tPlayer.Bakers); i++ {
		tPlayer.Bakers[i].Points -= bid
	}

	// Add player to current tournament
	tournament.Players[tPlayer.ID] = tPlayer
}

func ResultTournament(ctx *fasthttp.RequestCtx) {

	res, err := ValidateResult(ctx)

	// Return error if request is not valid
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetConnectionClose()
		return
	}

	tournament, err := db.Tournaments.Get(*res.ID)

	// Return error if tournament not exists
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.SetConnectionClose()
		return
	}

	// Calculate prize
	for _, winner := range *res.Winners {
		if player, ok := tournament.Players[*winner.ID]; ok {
			prize := *winner.Prize / float64(len(player.Bakers))
			for i := 0; i < len(player.Bakers); i++ {
				player.Bakers[i].Points += prize
			}
		}
	}

	// Set end flag for current tournament
	tournament.Ended = true
}
