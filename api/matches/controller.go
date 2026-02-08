package matches

import (
	"dotapro-lambda-api/errs"
	"dotapro-lambda-api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	model *Model
}

func NewController(model *Model) *Controller {
	return &Controller{model}
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
