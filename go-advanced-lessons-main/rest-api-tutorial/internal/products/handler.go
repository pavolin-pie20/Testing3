package products

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
	productsURL = "/products"
	productURL  = "/products/:uuid"
)

type handler struct {
	logger     *logging.Logger
	repository ProductRepository
}

func NewHandler(repository ProductRepository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, productsURL, apperror.Middleware(h.GetList))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	all, err := h.repository.FindAllProduct(context.TODO())
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
