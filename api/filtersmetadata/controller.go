package filtersmetadata

import (
	"context"
	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/utils"
	"net/http"
)

type Controller struct {
	model *Model
}

func NewController(model *Model) *Controller {
	return &Controller{model}
}

func (c *Controller) SearchTeams(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), constants.ShortRequestTimeout)
	defer cancel()

	query, err := utils.ParseRequiredStringParam(r.URL.Query(), "q")
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	teams, err := c.model.SearchTeams(ctx, query)
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

	utils.WriteResponse(w, teams, http.StatusOK)
}

func (c *Controller) SearchLeagues(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), constants.ShortRequestTimeout)
	defer cancel()

	query, err := utils.ParseRequiredStringParam(r.URL.Query(), "q")
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	leagues, err := c.model.SearchLeagues(ctx, query)
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

	utils.WriteResponse(w, leagues, http.StatusOK)
}

func (c *Controller) GetTeamName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), constants.ShortRequestTimeout)
	defer cancel()

	id, err := utils.ParseRequiredInt64Param(r.URL.Query(), "id")
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	name, err := c.model.GetTeamName(ctx, id)
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

	utils.WriteResponse(w, name, http.StatusOK)
}

func (c *Controller) GetLeagueName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), constants.ShortRequestTimeout)
	defer cancel()

	id, err := utils.ParseRequiredInt64Param(r.URL.Query(), "id")
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	name, err := c.model.GetLeagueName(ctx, id)
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

	utils.WriteResponse(w, name, http.StatusOK)
}
