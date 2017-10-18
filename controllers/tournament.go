package controllers

import (
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/db"
	"github.com/qiangxue/fasthttp-routing"
)

type ResultTournamentStruct struct {
	ID      *string   `json:"tournamentId"`
	Winners *[]Winner `json:"winners"`
}

type Winner struct {
	ID    *string  `json:"playerId"`
	Prize *float64 `json:"prize"`
}

const (
	QueryTournamentId = "tournamentId"
	QueryBackerId     = "backerId"
	QueryDeposit      = "deposit"
)

func AnnounceTournament(ctx *routing.Context) (err error) {
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

	// Check if "tournamentId" and "deposit" exists in query string
	if !query.Has(QueryTournamentId) || !query.Has(QueryDeposit) {
		respError = fasthttp.StatusBadRequest
		return
	}

	id := string(query.Peek(QueryTournamentId))
	deposit, err := query.GetUfloat(QueryDeposit)

	// Check if "deposit" is float and is non-negative number
	if err != nil || deposit < 0 {
		respError = fasthttp.StatusBadRequest
		return
	}

	tournament, err := db.Tournaments.Get(id, &conn)

	// Return error if tournament exists
	if tournament.ID != "" {
		respError = fasthttp.StatusNotAcceptable
		return
	}

	// Create tournament
	tournament = db.Tournament{ID: id, Deposit: deposit}
	tournament.Players = make(map[string]*db.TournamentPlayer)

	// Try save tournament
	err = tournament.Save(&conn)

	// Return error if tournament not saved
	if err != nil {
		respError = fasthttp.StatusBadRequest
		return
	}
	return
}

func JoinTournament(ctx *routing.Context) (e error) {
	conn := db.DB.Conn()
	query := ctx.QueryArgs()
	respError := fasthttp.StatusOK

	defer func() {
		if respError != fasthttp.StatusOK {
			conn.Rollback()
			ctx.SetStatusCode(respError)
			ctx.SetConnectionClose()
		}
		conn.Close()
	}()

	// Check if "tournamentId" and "playerId" exists in query string
	if !query.Has(QueryTournamentId) || !ctx.QueryArgs().Has(QueryPlayerId) {
		respError = fasthttp.StatusBadRequest
		return
	}

	tournamentId := string(ctx.QueryArgs().Peek(QueryTournamentId))
	playerId := string(ctx.QueryArgs().Peek(QueryPlayerId))
	backerIds := ctx.QueryArgs().PeekMulti(QueryBackerId)

	conn.BeginTransaction()

	tournament, err := db.Tournaments.Get(tournamentId, &conn)

	// Return error if tournament not exists
	if err != nil {
		respError = fasthttp.StatusNotFound
		return
	}

	// Return error if tournament ended
	if tournament.Ended {
		respError = fasthttp.StatusForbidden
		return
	}

	// Return error if player exists in current tournament
	if _, ok := tournament.Players[playerId]; ok {
		respError = fasthttp.StatusNotAcceptable
		return
	}

	// Calculate bid
	bid := tournament.Deposit / float64(len(backerIds)+1)

	player, err := db.Players.Get(playerId, &conn)

	// Return error if player not exists
	if err != nil {
		respError = fasthttp.StatusNotFound
		return
	}

	if player.Points-bid < 0 {
		respError = fasthttp.StatusPaymentRequired
		return
	}

	player.Points -= bid
	player.Save(&conn)

	// Create tournament player
	tournamentPlayer := &db.TournamentPlayer{ID: player.ID}

	tournamentPlayer.Bakers = make([]db.Player, len(backerIds)+1)
	tournamentPlayer.Bakers[0] = player

	// Find and add backers to player
	if len(backerIds) > 0 {
		for i := 0; i < len(backerIds); i++ {
			backer, err := db.Players.Get(string(backerIds[i]), &conn)
			// Return error if backer not exists
			if err != nil {
				respError = fasthttp.StatusNotFound
				return
			}

			// Check player balance
			if backer.Points-bid < 0 {
				respError = fasthttp.StatusPaymentRequired
				return
			}

			// Create bid
			backer.Points -= bid
			backer.Save(&conn)

			tournamentPlayer.Bakers[i+1] = backer
		}
	}

	// Add player to current tournament
	tournament.Players[tournamentPlayer.ID] = tournamentPlayer
	tournament.Save(&conn)
	conn.Commit()
	return
}

func ResultTournament(ctx *routing.Context) (err error) {
	conn := db.DB.Conn()
	respError := fasthttp.StatusOK

	defer func() {
		if respError != fasthttp.StatusOK {
			conn.Rollback()
			ctx.SetStatusCode(respError)
			ctx.SetConnectionClose()
		}
		conn.Close()
	}()

	res, err := ValidateResult(ctx)

	// Return error if request is not valid
	if err != nil {
		respError = fasthttp.StatusBadRequest
		return
	}

	conn.BeginTransaction()

	tournament, err := db.Tournaments.Get(*res.ID, &conn)

	// Return error if tournament not exists
	if err != nil {
		respError = fasthttp.StatusNotFound
		return
	}

	// Calculate prize
	for _, winner := range *res.Winners {
		if player, ok := tournament.Players[*winner.ID]; ok {
			prize := *winner.Prize / float64(len(player.Bakers))
			for i := 0; i < len(player.Bakers); i++ {

				player.Bakers[i].Points += prize
				player.Bakers[i].Save(&conn)
			}
		}
	}

	// Set end flag for current tournament
	tournament.Ended = true
	tournament.Save(&conn)
	conn.Commit()
	return
}
