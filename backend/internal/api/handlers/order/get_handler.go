package order

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/aliskhannn/order-service/internal/api/handlers"
	orderrepo "github.com/aliskhannn/order-service/internal/repository/order"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/aliskhannn/order-service/internal/model"
)

//go:generate mockgen -source=get_handler.go -destination=../../mocks/api/order/mock_order_service.go -package=mocks
type orderService interface {
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
}

type GetHandler struct {
	logger       *zap.Logger
	orderService orderService
}

func NewGetHandler(l *zap.Logger, s orderService) *GetHandler {
	return &GetHandler{
		logger:       l,
		orderService: s,
	}
}

func (h *GetHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderStr)
	if err != nil {
		h.logger.Error("failed to parse order id", zap.String("orderID", orderStr), zap.Error(err))
		handlers.WriteError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	if orderID == uuid.Nil {
		h.logger.Error("order id is empty", zap.String("orderID", orderStr))
		handlers.WriteError(w, http.StatusBadRequest, "order ID is required")
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, orderrepo.ErrOrderNotFound) {
			h.logger.Info("order not found", zap.String("orderID", orderID.String()))
			handlers.WriteError(w, http.StatusNotFound, "order not found")
			return
		}

		h.logger.Error("unexpected error getting order", zap.String("orderID", orderID.String()), zap.Error(err))
		handlers.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.logger.Info("order retrieved", zap.String("orderID", orderID.String()))

	handlers.WriteJSON(w, http.StatusOK, order)
}
