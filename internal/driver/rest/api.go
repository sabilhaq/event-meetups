package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gopkg.in/validator.v2"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/play"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
)

const (
	publicDir = "/dist"
	indexFile = "index.html"
)

type API struct {
	serviceName   string
	playService   play.Service
	battleService battle.Service
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(otelchi.Middleware(
		a.serviceName,
		otelchi.WithChiRoutes(r),
		otelchi.WithPropagators(otel.GetTextMapPropagator()),
		otelchi.WithTracerProvider(otel.GetTracerProvider()),
		otelchi.WithRequestMethodInSpanName(true),
	))
	r.Get("/partners", a.serveGetAvailablePartners)
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
	// serve the frontend in SPA mode
	r.NotFound(a.serveWebFrontend)
	return r
}

func (a *API) serveWebFrontend(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Clean(r.URL.Path)
	if fileName != "index.html" && !strings.Contains(fileName, "assets") {
		fileName = "assets" + fileName
	}
	p := filepath.Join(publicDir, fileName)

	if info, err := os.Stat(p); err != nil {
		http.ServeFile(w, r, filepath.Join(publicDir, indexFile))
		return
	} else if info.IsDir() {
		http.ServeFile(w, r, filepath.Join(publicDir, indexFile))
		return
	}

	http.ServeFile(w, r, p)
}

func (a *API) serveGetAvailablePartners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveGetAvailablePartners: /partners")
	defer span.End()

	partners, err := a.playService.GetAvailablePartners(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"partners": partners,
	}))
}

func (a *API) serveNewGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveNewGame: POST /games/")
	defer span.End()

	var rb newGameReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	err = rb.Validate()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		render.Render(w, r, NewErrorResp(err))
		return
	}
	game, err := a.playService.NewGame(ctx, rb.PlayerName, rb.PartnerID)
	if err != nil {
		if errors.Is(err, play.ErrPartnerNotFound) {
			err = NewPartnerNotFoundError()
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetGameDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveGetGameDetails: GET /games/{game_id}")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

	game, err := a.playService.GetGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, play.ErrGameNotFound) {
			err = NewGameNotFoundError()
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(game))
}

func (a *API) serveGetScenario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveGetScenario: GET /games/{game_id}/scenario")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

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

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveStartBattle: PUT /games/{game_id}/battle")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

	bt, err := a.battleService.StartBattle(ctx, gameID)
	if err != nil {
		switch err {
		case battle.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battle.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		case battle.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveGetBattleInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveGetBattleInfo: GET /games/{game_id}/battle")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

	bt, err := a.battleService.GetBattle(ctx, gameID)
	if err != nil {
		switch err {
		case battle.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battle.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		case battle.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveDecideTurn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveDecideTurn: PUT /games/{game_id}/battle/turn")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

	bt, err := a.battleService.DecideTurn(ctx, gameID)
	if err != nil {
		switch err {
		case battle.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battle.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		case battle.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveAttack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveAttack: PUT /games/{game_id}/battle/attack")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

	bt, err := a.battleService.Attack(ctx, gameID)
	if err != nil {
		switch err {
		case battle.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battle.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		case battle.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

func (a *API) serveSurrender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// tracing
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "serveSurrender: PUT /games/{game_id}/battle/surrender")
	defer span.End()

	gameID := chi.URLParam(r, "game_id")
	span.SetAttributes(attribute.Key("game-id").String(gameID))

	bt, err := a.battleService.Surrender(ctx, gameID)
	if err != nil {
		switch err {
		case battle.ErrGameNotFound:
			err = NewGameNotFoundError()
		case battle.ErrBattleNotFound:
			err = NewBattleNotFoundError()
		case battle.ErrInvalidBattleState:
			err = NewInvalidBattleStateError()
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(bt))
}

type APIConfig struct {
	PlayingService play.Service   `validate:"nonnil"`
	BattleService  battle.Service `validate:"nonnil"`
	ServiceName    string         `validate:"nonzero"`
}

func (c APIConfig) Validate() error {
	return validator.Validate(c)
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	a := &API{
		playService:   cfg.PlayingService,
		battleService: cfg.BattleService,
		serviceName:   cfg.ServiceName,
	}
	return a, nil
}
