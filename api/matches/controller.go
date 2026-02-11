package matches

import (
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

	if pageStr := params.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid page: %v", err.Error()), http.StatusBadRequest)
			return
		}
		filter.Page = page
	}

	matches, paginationData, err := c.model.GetMany(r.Context(), filter)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(matches) == 0 {
		utils.WriteError(w, errs.NOT_FOUND.Error(), http.StatusNotFound)
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
