package main

import (
	"log"
	"runtime"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/controllers"
)

func main() {
	runtime.GOMAXPROCS(1)

	router := routing.New()

	// Player routes
	router.Get("/fund", controllers.PlayerFund)
	router.Get("/take", controllers.PlayerTake)
	router.Get("/balance", controllers.PlayerBalance)

	// Tournament routes
	router.Get("/announceTournament", controllers.AnnounceTournament)
	router.Get("/joinTournament", controllers.JoinTournament)
	router.Post("/resultTournament", controllers.ResultTournament)

	// Reset route
	router.Get("/reset", controllers.Reset)

	// Serve
	log.Fatal(fasthttp.ListenAndServe(":80", router.HandleRequest))
}
