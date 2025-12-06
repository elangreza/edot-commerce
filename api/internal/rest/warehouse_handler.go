package rest

import (
	"context"
	"encoding/json"
	"net/http"

	errs "github.com/elangreza/edot-commerce/api/internal/error"
	"github.com/elangreza/edot-commerce/api/internal/params"
	"github.com/go-chi/chi/v5"
)

type (
	WarehouseService interface {
		SetWarehouseStatus(ctx context.Context, req params.SetWarehouseStatusRequest) error
	}

	WarehouseHandler struct {
		svc WarehouseService
	}
)

func NewWarehouseHandler(
	publicRoute chi.Router,
	authService AuthService,
	svc WarehouseService,
) {

	authMiddleware := AuthMiddleware{
		svc: authService,
	}

	oh := WarehouseHandler{
		svc: svc,
	}

	publicRoute.Group(func(r chi.Router) {
		r.Use(authMiddleware.MustAuthMiddleware())
		r.Post("/warehouse/status", oh.CreateWarehouse())
	})
}

func (oh *WarehouseHandler) CreateWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := params.SetWarehouseStatusRequest{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, errs.ValidationError{Message: err.Error()})
			return
		}

		ctx := r.Context()

		err := oh.svc.SetWarehouseStatus(ctx, body)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		sendSuccessResponse(w, http.StatusCreated, "ok")
	}
}
