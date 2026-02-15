package series

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
	Series     []types.SeriesSummary `json:"series"`
	Pagination types.PaginationData  `json:"pagination"`
}

func NewController(model *Model) *Controller {
	return &Controller{model}
}

func (c *Controller) GetMany(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	filter := types.GetSeriesFilter{}
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

	filter.Sort = params.Get("sort")

	if limitStr := params.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid limit: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.Limit = limit
	}

	if pageStr := params.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid page: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.Page = page
	}

	series, paginationData, err := c.model.GetMany(ctx, filter)
	if err != nil {
		if err == context.Canceled {
			log.Debug().Msg("Request canceled by client")
			return
		}
		if err == context.DeadlineExceeded {
			log.Debug().Msg("Deadline exceeded ")
			utils.WriteError(w, context.DeadlineExceeded.Error(), http.StatusGatewayTimeout)
			return
		}
		log.Err(err).Send()
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetManyResp{Series: series, Pagination: paginationData}
	utils.WriteResponse(w, resp, http.StatusOK)
}

func (c *Controller) GetOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("missing or invalid id: %v", err.Error()), http.StatusBadRequest)
		return
	}
	series, err := c.model.GetOne(r.Context(), id)
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
	utils.WriteResponse(w, series, http.StatusOK)
}
