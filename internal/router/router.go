package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ibraheemdev/poller/config"
	"github.com/ibraheemdev/poller/internal/polls"
	"github.com/ibraheemdev/poller/pkg/middleware"
	"github.com/julienschmidt/httprouter"
)

// ListenAndServe :
func ListenAndServe() {
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port),
			initRoutes()))
}

func initRoutes() *httprouter.Router {
	r := httprouter.New()
	r.POST("/polls", middleware.Cors(polls.Create()))
	r.GET("/polls/:id", middleware.Cors(polls.Show()))
	r.PUT("/polls/:id", middleware.Cors(polls.Update()))
	return r
}
