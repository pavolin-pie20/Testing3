package customers

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"restapi-lesson/internal/apperror"
	"restapi-lesson/internal/handlers"
	"restapi-lesson/pkg/logging"
)

const (
	columnsURL = "/columns"
	columnURL  = "/columns/:uuid"
)

type handler struct {
	logger     *logging.Logger
	repository CustomerRepository
}

func NewHandler(repository CustomerRepository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, columnsURL, apperror.Middleware(h.GetList))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	all, err := h.repository.FindAllCustomers(context.TODO())
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}
