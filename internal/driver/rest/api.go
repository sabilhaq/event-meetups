package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/event"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/meetup"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/session"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/venue"
)

type APIConfig struct {
	PlayingService play.Service    `validate:"nonnil"`
	BattleService  battle.Service  `validate:"nonnil"`
	SessionService session.Service `validate:"nonnil"`
	EventService   event.Service   `validate:"nonnil"`
	VenueService   venue.Service   `validate:"nonnil"`
	MeetupService  meetup.Service  `validate:"nonnil"`
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
		sessionService: cfg.SessionService,
		eventService:   cfg.EventService,
		venueService:   cfg.VenueService,
		meetupService:  cfg.MeetupService,
		isWebEnabled:   cfg.IsWebEnabled,
	}
	return a, nil
}

type API struct {
	playService    play.Service
	battleService  battle.Service
	sessionService session.Service
	eventService   event.Service
	venueService   venue.Service
	meetupService  meetup.Service
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

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		r.Get("/events", a.serveGetEvents)
		r.Route("/venues", func(r chi.Router) {
			r.Get("/", a.serveGetVenues)
			r.Get("/{venue_id}", a.serveGetVenue)
		})

		r.Route("/meetups", func(r chi.Router) {
			r.Post("/", a.serveCreateMeetup)
			r.Get("/", a.serveGetMeetups)
			r.Route("/{meetup_id}", func(r chi.Router) {
				r.Get("/", a.serveGetMeetup)
				r.Put("/", a.serveUpdateMeetup)
				r.Delete("/", a.serveCancelMeetup)
			})
		})

		r.Route("/incoming-meetups", func(r chi.Router) {
			r.Get("/", a.serveGetIncomingMeetups)
			r.Route("/{meetup_id}", func(r chi.Router) {
				r.Put("/", a.serveJoinMeetup)
				r.Delete("/", a.serveLeaveMeetup)
			})
		})
	})

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

func (a *API) serveGetVenues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filter := entity.GetVenueFilter{}

	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr != "" {
		eventID, err := strconv.Atoi(eventIDStr)
		if err != nil {
			render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		}
		filter.EventID = &eventID
	}

	meetupStartTSStr := r.URL.Query().Get("meetup_start_ts")
	if meetupStartTSStr != "" {
		meetupStartTS, err := strconv.ParseInt(meetupStartTSStr, 10, 64)
		if err != nil {
			render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		}
		meetupStartTSHHMM := time.Unix(meetupStartTS, 0).Format("15:04")
		filter.MeetupStartTS = &meetupStartTSHHMM
	}

	meetupEndTSStr := r.URL.Query().Get("meetup_end_ts")
	if meetupEndTSStr != "" {
		meetupEndTS, err := strconv.ParseInt(meetupEndTSStr, 10, 64)
		if err != nil {
			render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		}
		meetupEndTSHHMM := time.Unix(meetupEndTS, 0).Format("15:04")
		filter.MeetupEndTS = &meetupEndTSHHMM
	}

	venues, err := a.venueService.GetVenues(ctx, filter)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(venues))
}

func (a *API) serveGetVenue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	venueID, err := strconv.Atoi(chi.URLParam(r, "venue_id"))
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	venue, err := a.venueService.GetVenue(ctx, venueID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(venue))
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

func (a *API) serveCreateMeetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())

	var rb createMeetupReqBody
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
	meetup, err := a.meetupService.CreateMeetup(ctx, entity.CreateMeetupRequest{
		Name:        rb.Name,
		VenueID:     rb.VenueID,
		EventID:     rb.EventID,
		StartTs:     rb.StartTs,
		EndTs:       rb.EndTs,
		MaxPersons:  rb.MaxPersons,
		OrganizerID: userID,
	})
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(meetup))
}

func (a *API) serveGetMeetups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filter := entity.GetMeetupFilter{}

	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr != "" {
		eventID, err := strconv.Atoi(eventIDStr)
		if err != nil {
			render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		}
		filter.EventID = &eventID
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		}
		filter.Limit = &limit
	}

	meetups, err := a.meetupService.GetMeetups(ctx, filter)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(map[string]interface{}{
		"meetups": meetups,
	}))
}

func (a *API) serveGetMeetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())

	meetupID, err := strconv.Atoi(chi.URLParam(r, "meetup_id"))
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	meetup, err := a.meetupService.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(meetup))
}

func (a *API) serveUpdateMeetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())

	meetupID, err := strconv.Atoi(chi.URLParam(r, "meetup_id"))
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	var rb updateMeetupReqBody
	err = json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	err = rb.Validate()
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	meetup, err := a.meetupService.UpdateMeetup(ctx, meetupID, entity.UpdateMeetupRequest{
		Name:       rb.Name,
		StartTs:    rb.StartTs,
		EndTs:      rb.EndTs,
		MaxPersons: rb.MaxPersons,
		UserID:     userID,
	})
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(meetup))
}

func (a *API) serveCancelMeetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())

	meetupID, err := strconv.Atoi(chi.URLParam(r, "meetup_id"))
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	cancelledReason := r.URL.Query().Get("cancelled_reason")
	if cancelledReason == "" {
		render.Render(w, r, NewErrorResp(NewCancelledReasonRequiredError()))
		return
	}

	meetup, err := a.meetupService.CancelMeetup(ctx, meetupID, userID, cancelledReason)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(meetup))
}

func (a *API) serveJoinMeetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())

	meetupID, err := strconv.Atoi(chi.URLParam(r, "meetup_id"))
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	meetup, err := a.meetupService.JoinMeetup(ctx, meetupID, userID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(meetup))
}

func (a *API) serveLeaveMeetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())

	meetupID, err := strconv.Atoi(chi.URLParam(r, "meetup_id"))
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	err = a.meetupService.LeaveMeetup(ctx, meetupID, userID)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}
	render.Render(w, r, NewSuccessResp(nil))
}

func (a *API) serveGetIncomingMeetups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserFromContext(r.Context())
	filter := entity.GetIncomingMeetupFilter{}
	filter.UserID = userID

	filter.Status = r.URL.Query().Get("status")
	if filter.Status == "" {
		filter.Status = "all"
	}
	if filter.Status != "all" && filter.Status != "open" && filter.Status != "cancelled" {
		render.Render(w, r, NewErrorResp(NewBadRequestError("status: invalid status")))
		return
	}

	eventIDs := r.URL.Query().Get("event_ids")
	if eventIDs != "" {
		filter.EventIDs = &eventIDs
	}

	venueIDs := r.URL.Query().Get("venue_ids")
	if venueIDs != "" {
		filter.VenueIDs = &venueIDs
	}

	meetups, err := a.meetupService.GetIncomingMeetups(ctx, entity.GetIncomingMeetupFilter{
		UserID:   filter.UserID,
		EventIDs: filter.EventIDs,
		VenueIDs: filter.VenueIDs,
		Status:   filter.Status,
	})
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(meetups))
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
	case venue.ErrVenueNotFound:
		err = NewNotFoundError()
	case meetup.ErrInvalidEvent:
		err = NewInvalidEventError()
	case meetup.ErrExceedVenueCapacity:
		err = NewExceedVenueCapacityError()
	case meetup.ErrVenueIsClosed:
		err = NewVenueIsClosedError()
	case meetup.ErrForbidden:
		err = NewForbiddenError()
	case meetup.ErrMaxPersonsLessThanJoinedPersons:
		err = NewMaxPersonsLessThanJoinedPersonsError()
	case meetup.ErrCancelledReasonRequred:
		err = NewCancelledReasonRequiredError()
	case meetup.ErrMeetupStarted:
		err = NewMeetupStartedError()
	case meetup.ErrMeetupFinished:
		err = NewMeetupFinishedError()
	case meetup.ErrMeetupCancelled:
		err = NewMeetupCancelledError()
	case meetup.ErrMeetupClosed:
		err = NewMeetupClosedError()
	case meetup.ErrMeetupOverlaps:
		err = NewMeetupOverlapsError()
	case meetup.ErrUserNotParticipant:
		err = NewUserNotParticipantError()
	default:
		err = NewInternalServerError(err.Error())
	}
	render.Render(w, r, NewErrorResp(err))
}
