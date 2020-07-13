package app

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"geogame/internal/locations"
	"geogame/internal/middleware"
	"geogame/internal/players"
	"geogame/pkg"
)

// type assert the main controller which extend the chi controller
var _ pkg.Controller = (*Controller)(nil)

type Controller struct {
	logger    *zap.Logger
	locations locations.Service
	players   players.Service
	jwtAuther middleware.JwtAuther
}

func NewController(logger *zap.Logger, locations locations.Service, players players.Service, jwtAuther middleware.JwtAuther) *Controller {
	return &Controller{
		logger:    logger,
		locations: locations,
		players:   players,
		jwtAuther: jwtAuther,
	}
}

// extend the chi controller init
func (c *Controller) Init(logger *zap.Logger) error {
	c.logger = logger
	return nil
}

// extend the chi controller
// Setup the chi routes
func (c *Controller) SetupRouter(router chi.Router) error {

	// Register admin endpoints
	router.Route("/admin/loc", func(r chi.Router) {
		r.Post("/create", c.CreateLocation)
		r.Get("/{id}", c.GetLocation)
		r.Put("/update", c.UpdateLocation)
		r.Delete("/{id}/delete", c.DeleteLocation)
	})

	// Register client endpoints
	router.Route("/client", func(r chi.Router) {
		r.Post("/register", c.Register)
		r.Post("/login", c.Login)
		r.With(middleware.IsClientAllowed(c.jwtAuther)).Post("/loc/send", c.SendLocation)
		r.With(middleware.IsClientAllowed(c.jwtAuther)).Put("/update-name", c.UpdateName)
		r.With(middleware.IsClientAllowed(c.jwtAuther)).Get("/loc/get", c.GetClientLocation)
	})

	return nil
}

// extend the chi controller init
func (c *Controller) Terminate() error {
	return nil
}
