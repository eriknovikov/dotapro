package matches

import (
	"dotapro-lambda-api/errs"
	"net/http"
)

type Controller struct {
	model *Model
}

func NewController(model *Model) *Controller {
	return &Controller{model}
}

func (c *Controller) GetOne(w http.ResponseWriter, r *http.Request) {
	http.Error(w, errs.UNIMPLEMENTED.Error(), http.StatusNotImplemented)
}
