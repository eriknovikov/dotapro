package matches

import (
	"context"
	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/errs"
	"dotapro-lambda-api/types"
	"dotapro-lambda-api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	ctx, cancel := context.WithTimeout(r.Context(), constants.DefaultRequestTimeout)
	defer cancel()
	
	filter := types.GetMatchesFilter{}
	params := r.URL.Query()

	// Parse optional filter parameters
	if leagueID, err := utils.ParseInt64Param(params, "league"); err != nil {
		utils.WriteParamError(w, "league_id", err)
		return
	} else {
		filter.LeagueID = leagueID
	}

	if teamID, err := utils.ParseInt64Param(params, "team"); err != nil {
		utils.WriteParamError(w, "team_id", err)
		return
	} else {
		filter.TeamID = teamID
	}

	if playerID, err := utils.ParseInt64Param(params, "player"); err != nil {
		utils.WriteParamError(w, "player_id", err)
		return
	} else {
		filter.PlayerID = playerID
	}

	if heroID, err := utils.ParseInt64Param(params, "hero"); err != nil {
		utils.WriteParamError(w, "hero_id", err)
		return
	} else {
		filter.HeroID = heroID
	}

	// Parse sort parameter
	filter.Sort = utils.ParseStringParam(params, "sort")

	// Parse limit parameter
	if limit, err := utils.ParseIntParam(params, "limit"); err != nil {
		utils.WriteParamError(w, "limit", err)
		return
	} else {
		filter.Limit = limit
	}

	// Parse cursor parameter
	if cursor, err := utils.ParseInt64Param(params, "c"); err != nil {
		utils.WriteParamError(w, "cursor", err)
		return
	} else {
		filter.Cursor = cursor
	}

	matches, paginationData, err := c.model.GetMany(ctx, filter)
	if err != nil {
		if err == context.Canceled {
			return
		}

		if err == context.DeadlineExceeded {
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
	id, err := utils.ParseRequiredInt64Param(r.URL.Query(), "id")
	if err != nil {
		// For URL path parameters, we need to use chi.URLParam
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			utils.WriteError(w, "missing id parameter", http.StatusBadRequest)
			return
		}
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			utils.WriteError(w, fmt.Sprintf("invalid id: %v", err), http.StatusBadRequest)
			return
		}
	}
	match, err := c.model.GetOne(r.Context(), id)
	if err != nil {
		// Handle context cancellation gracefully - client aborted the request
		if err == context.Canceled || err == context.DeadlineExceeded {
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
