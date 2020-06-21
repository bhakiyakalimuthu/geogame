package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"geogame/internal/locations"
	"geogame/internal/middleware"
	"geogame/internal/players"
	"geogame/pkg"
)

const (
	HTTPContentType     string = "Content-Type"
	HTTPApplicationJSON string = "application/json"
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

// admin endpoints
func (c *Controller) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var payload locations.Location
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.locations.Create(r.Context(), payload); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, nil)
}

func (c *Controller) GetLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := c.locations.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, res)
}

func (c *Controller) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var payload locations.Location
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.locations.Update(r.Context(), payload); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, nil)
}

func (c *Controller) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := c.locations.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, nil)
}

// client endpoints

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var payload players.RegisterPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeResponse(w, http.StatusBadRequest, err)
		return
	}
	if payload.Email == "" {
		c.logger.Error("Register: failed to register", zap.Error(errors.New("empty email id")))
		writeError(w, http.StatusBadRequest, errors.New("empty email id"))
		return
	}
	if payload.Password == "" {
		c.logger.Error("Register: failed to register", zap.Error(errors.New("empty password")))
		writeError(w, http.StatusBadRequest, errors.New("empty password id"))
		return
	}
	if payload.Name == "" {
		c.logger.Error("Register: failed to register", zap.Error(errors.New("empty name")))
		writeError(w, http.StatusBadRequest, errors.New("empty name"))
		return
	}
	if err := c.players.Register(r.Context(), payload); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, nil)

}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var payload players.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if payload.Email == "" {
		c.logger.Error("Login: failed to register", zap.Error(errors.New("empty email id")))
		writeError(w, http.StatusBadRequest, errors.New("empty email id"))
		return
	}
	if payload.Password == "" {
		c.logger.Error("Login: failed to register", zap.Error(errors.New("empty password")))
		writeError(w, http.StatusBadRequest, errors.New("empty password id"))
		return
	}
	res, err := c.players.Login(r.Context(), payload)
	if err != nil {

		writeError(w, http.StatusInternalServerError, err)
		return
	}
	// c.logger.Info("res",zap.Any("jwt",res))
	writeResponse(w, http.StatusOK, *res)
}

func (c *Controller) SendLocation(w http.ResponseWriter, r *http.Request) {
	token, err := extractTokenFromContext(r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var p locations.Location
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.players.UpdateLocation(r.Context(), p, token.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, nil)
}

func (c *Controller) UpdateName(w http.ResponseWriter, r *http.Request) {
	token, err := extractTokenFromContext(r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var p players.UpdatePayload
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if p.Name == "" {
		writeError(w, http.StatusBadRequest, errors.New("empty name"))
		return
	}

	if err := c.players.UpdateName(r.Context(), p, token.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, nil)
}

func (c *Controller) GetClientLocation(w http.ResponseWriter, r *http.Request) {

	token, err := extractTokenFromContext(r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	res, err := c.players.GetLocation(r.Context(), token.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeResponse(w, http.StatusOK, res)
}

func extractTokenFromContext(r *http.Request) (*middleware.AccessToken, error) {
	token, ok := r.Context().Value("AccessToken").(*middleware.AccessToken)
	if !ok || token == nil {
		return nil, errors.New("failed to extract access token from context")
	}
	return token, nil
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set(HTTPContentType, HTTPApplicationJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func writeError(w http.ResponseWriter, statusCode int, httpError error) {

	w.Header().Set(HTTPContentType, HTTPApplicationJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&ErrorResponse{Error: httpError.Error()}); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

type ErrorResponse struct {
	Error string
}
