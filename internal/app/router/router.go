package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ibraheemdev/agilely/config"
	"github.com/ibraheemdev/agilely/internal/app/polls"
	"github.com/ibraheemdev/agilely/pkg/middleware"
	"github.com/julienschmidt/httprouter"
)

// ListenAndServe :
func ListenAndServe(r *httprouter.Router) {
	initRoutes(r)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port),
			r))
}

func initRoutes(r *httprouter.Router) {
	r.POST("/polls", middleware.Cors(polls.Create()))
	r.GET("/polls/:id", middleware.Cors(polls.Show()))
	r.PUT("/polls/:id", middleware.Cors(polls.Update()))
}
