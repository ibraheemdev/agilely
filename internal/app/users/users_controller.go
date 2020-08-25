package users

import (
	"github.com/ibraheemdev/agilely/internal/app/engine"
)

// Users controller
type Users struct {
	*engine.Engine
}

// NewController : Returns a new users controller
func NewController(e *engine.Engine) *Users {
	return &Users{Engine: e}
}
