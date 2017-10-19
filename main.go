package main

import (
	"log"
	"runtime"
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/controllers"
)

const (
	// Player routes
	pathFund    = "/fund"
	pathTake    = "/take"
	pathBalance = "/balance"

	// Tournament routes
	pathAnnounceTournament = "/announceTournament"
	pathJoinTournament     = "/joinTournament"
	pathResultTournament   = "/resultTournament"

	// Reset route
	pathReset = "/reset"
)

func HandleRequest(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	switch path {
	case pathFund:
		controllers.PlayerFund(ctx)
	case pathTake:
		controllers.PlayerTake(ctx)
	case pathBalance:
		controllers.PlayerBalance(ctx)
	case pathAnnounceTournament:
		controllers.AnnounceTournament(ctx)
	case pathJoinTournament:
		controllers.JoinTournament(ctx)
	case pathResultTournament:
		controllers.ResultTournament(ctx)
	case pathReset:
		controllers.Reset(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Serve
	log.Fatal(fasthttp.ListenAndServe(":80", HandleRequest))
}
