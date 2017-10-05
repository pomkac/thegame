package main

import (
	"github.com/buaazp/fasthttprouter"
	"log"
	"github.com/valyala/fasthttp"
	"github.com/pomkac/thegame/controllers"
)

func main() {
	router := fasthttprouter.New()

	// Player routes
	router.GET("/fund", controllers.PlayerFund)
	router.GET("/take", controllers.PlayerTake)
	router.GET("/balance", controllers.PlayerBalance)

	// Tournament routes
	router.GET("/announceTournament", controllers.AnnounceTournament)
	router.GET("/joinTournament", controllers.JoinTournament)
	router.POST("/resultTournament", controllers.ResultTournament)

	// Reset route
	router.GET("/reset", controllers.Reset)

	// Serve
	log.Fatal(fasthttp.ListenAndServe(":80", router.Handler))
}
