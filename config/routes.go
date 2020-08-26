package config

import (
	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/internal/app/polls"
	"github.com/ibraheemdev/agilely/internal/app/users"
)

// Routes :
func Routes(e *engine.Engine) error {
	// polls router
	p := polls.NewController(e)
	p.Core.Router.POST("/polls", p.Core.ErrorHandler.Wrap(p.Create))
	p.Core.Router.GET("/polls/:id", p.Core.ErrorHandler.Wrap(p.Show))
	p.Core.Router.PUT("/polls/:id", p.Core.ErrorHandler.Wrap(p.Update))

	// users routes
	u := users.NewController(e)
	if err := u.Init(); err != nil {
		return err
	}

	u.Core.Router.GET("/login", u.Core.ErrorHandler.Wrap(u.LoginGet))
	u.Core.Router.POST("/login", u.Core.ErrorHandler.Wrap(u.LoginPost))

	u.Core.Router.DELETE("/logout", u.Core.ErrorHandler.Wrap(u.Logout))

	u.Core.Router.GET("/confirm", u.Core.ErrorHandler.Wrap(u.GetConfirm))

	u.Core.Router.GET("/recover", u.Core.ErrorHandler.Wrap(u.StartGetRecover))
	u.Core.Router.POST("/recover", u.Core.ErrorHandler.Wrap(u.StartPostRecover))
	u.Core.Router.GET("/recover/end", u.Core.ErrorHandler.Wrap(u.EndGetRecover))
	u.Core.Router.POST("/recover/end", u.Core.ErrorHandler.Wrap(u.EndPostRecover))

	u.Core.Router.GET("/register", u.Core.ErrorHandler.Wrap(u.GetRegister))
	u.Core.Router.POST("/register", u.Core.ErrorHandler.Wrap(u.PostRegister))

	return nil
}
