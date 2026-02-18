package matches

import (
	"context"
	"dotapro-lambda-api/errs"
	"dotapro-lambda-api/types"
	"dotapro-lambda-api/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	model *Model
}

type GetManyResp struct {
	Matches    []types.MatchSummary `json:"matches"`
	Pagination types.PaginationData `json:"pagination"`
}

func NewController(model *Model) *Controller {
	return &Controller{model}
}

func (c *Controller) GetMany(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	filter := types.GetMatchesFilter{}
	params := r.URL.Query()

	if leagueIDStr := params.Get("league"); leagueIDStr != "" {
		leagueID, err := strconv.ParseInt(leagueIDStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid league_id: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.LeagueID = &leagueID
	}

	if teamIDStr := params.Get("team"); teamIDStr != "" {
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid team_id: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.TeamID = &teamID
	}

	if playerIDStr := params.Get("player"); playerIDStr != "" {
		playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid player_id: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.PlayerID = &playerID
	}

	if heroIDStr := params.Get("hero"); heroIDStr != "" {
		heroID, err := strconv.ParseInt(heroIDStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid hero_id: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.HeroID = &heroID
	}

	filter.Sort = params.Get("sort")

	if limitStr := params.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid limit: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.Limit = limit
	}

	if cursorStr := params.Get("c"); cursorStr != "" {
		cursor, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid cursor: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.Cursor = &cursor
	}

	matches, paginationData, err := c.model.GetMany(ctx, filter)
	if err != nil {
		if err == context.Canceled {
			return
		}

		if err == context.DeadlineExceeded {
			log.Debug().Msg("Deadline exceeded ")
			utils.WriteError(w, context.DeadlineExceeded.Error(), http.StatusGatewayTimeout)
			return
		}
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetManyResp{Matches: matches, Pagination: paginationData}
	utils.WriteResponse(w, resp, http.StatusOK)
}

func (c *Controller) GetOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("missing or invalid id: %v", err.Error()), http.StatusBadRequest)
		return
	}
	match, err := c.model.GetOne(r.Context(), id)
	if err != nil {
		// Handle context cancellation gracefully - client aborted the request
		if err == context.Canceled || err == context.DeadlineExceeded {
			log.Debug().Msg("Request canceled by client")
			return
		}
		switch err {
		case errs.NOT_FOUND:
			utils.WriteError(w, err.Error(), http.StatusNotFound)
		default:
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	utils.WriteResponse(w, match, http.StatusOK)
}
