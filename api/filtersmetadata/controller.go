package filtersmetadata

import (
	"context"
	"dotapro-lambda-api/utils"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Controller struct {
	model *Model
}

func NewController(model *Model) *Controller {
	return &Controller{model}
}

func (c *Controller) SearchTeams(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteError(w, "missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	teams, err := c.model.SearchTeams(ctx, query)
	if err != nil {
		if err == context.Canceled {
			return
		}
		if err == context.DeadlineExceeded {
			log.Debug().Msg("Deadline exceeded")
			utils.WriteError(w, context.DeadlineExceeded.Error(), http.StatusGatewayTimeout)
			return
		}
		log.Err(err).Send()
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, teams, http.StatusOK)
}

func (c *Controller) SearchLeagues(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteError(w, "missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	leagues, err := c.model.SearchLeagues(ctx, query)
	if err != nil {
		if err == context.Canceled {
			return
		}
		if err == context.DeadlineExceeded {
			log.Debug().Msg("Deadline exceeded")
			utils.WriteError(w, context.DeadlineExceeded.Error(), http.StatusGatewayTimeout)
			return
		}
		log.Err(err).Send()
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, leagues, http.StatusOK)
}
