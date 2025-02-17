package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/event"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/session"
)

type APIConfig struct {
	PlayingService play.Service    `validate:"nonnil"`
	BattleService  battle.Service  `validate:"nonnil"`
	EventService   event.Service   `validate:"nonnil"`
	SessionService session.Service `validate:"nonnil"`
	IsWebEnabled   bool
}

func (c APIConfig) Validate() error {
	return validator.Validate(c)
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	a := &API{
		playService:    cfg.PlayingService,
		battleService:  cfg.BattleService,
		eventService:   cfg.EventService,
		sessionService: cfg.SessionService,
		isWebEnabled:   cfg.IsWebEnabled,
	}
	return a, nil
}

type API struct {
	playService    play.Service
	battleService  battle.Service
	eventService   event.Service
	sessionService session.Service
	isWebEnabled   bool
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	if a.isWebEnabled {
		// by default route everything to the web client
		r.NotFound(a.serveWebClient)
	}

	r.Get("/health", a.serveHealthCheck)
	r.Get("/partners", a.serveGetAvailablePartners)
	r.Post("/session", a.serveCreateSession)
	r.Get("/events", a.serveGetEvents)
	r.Route("/games", func(r chi.Router) {
		r.Post("/", a.serveNewGame)
		r.Route("/{game_id}", func(r chi.Router) {
			r.Get("/", a.serveGetGameDetails)
			r.Get("/scenario", a.serveGetScenario)
			r.Route("/battle", func(r chi.Router) {
				r.Put("/", a.serveStartBattle)
				r.Get("/", a.serveGetBattleInfo)
				r.Put("/turn", a.serveDecideTurn)
				r.Put("/attack", a.serveAttack)
				r.Put("/surrender", a.serveSurrender)
			})
		})
	})

	return r
}

const (
	publicDir  = "./client"
	indexFile  = "index.html"
	assetsPath = "assets"
)

func (a *API) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *API) serveWebClient(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Clean(r.URL.Path)
	if fileName != indexFile && !strings.Contains(fileName, assetsPath) {
		fileName = assetsPath + fileName
	}
	p := filepath.Join(publicDir, fileName)

	if info, err := os.Stat(p); err != nil || info.IsDir() {
		http.ServeFile(w, r, filepath.Join(publicDir, indexFile))
		return
	}

	http.ServeFile(w, r, p)
}

func (a *API) serveGetAvailablePartners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	partners, err := a.playService.GetAvailablePartners(ctx)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"partners": partners,
	}))
}

func (a *API) serveGetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := a.eventService.GetEvents(ctx)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(events))
}

func (a *API) serveCreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rb newSessionReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	err = rb.Validate()
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	session, err := a.sessionService.CreateSession(ctx, rb.Username, rb.Password)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(session))
}

func (a *API) serveNewGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rb newGameReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	err = rb.Validate()
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	game, err := a.playService.NewGame(ctx, rb.PlayerName, rb.PartnerID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetGameDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	game, err := a.playService.GetGame(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetScenario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	game, err := a.playService.GetGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, play.ErrGameNotFound) {
			err = NewGameNotFoundError()
		}
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"scenario": game.Scenario,
	}))
}

func (a *API) serveStartBattle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.StartBattle(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveGetBattleInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.GetBattle(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveDecideTurn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.DecideTurn(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveAttack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.Attack(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveSurrender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID := chi.URLParam(r, "game_id")
	bt, err := a.battleService.Surrender(ctx, gameID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case battle.ErrGameNotFound:
		err = NewGameNotFoundError()
	case battle.ErrBattleNotFound:
		err = NewBattleNotFoundError()
	case battle.ErrInvalidBattleState:
		err = NewInvalidBattleStateError()
	case play.ErrGameNotFound:
		err = NewGameNotFoundError()
	case play.ErrPartnerNotFound:
		err = NewPartnerNotFoundError()
	case session.ErrInvalidCreds:
		err = NewSessionInvalidCredsError()
	default:
		err = NewInternalServerError(err.Error())
	}
	render.Render(w, r, NewErrorResp(err))
}
