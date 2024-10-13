package handler

import (
	"encoding/json"
	"kitchen/services/common/errors"
	"kitchen/services/common/genproto/orders"
	"kitchen/services/common/helpers"
	"kitchen/services/orders/entity"
	"kitchen/services/orders/service"
	"net/http"
)

type OrdersHTTPHandler interface {
	CreateOrder(http.ResponseWriter, *http.Request)
}

type ordersHTTPHandler struct {
	service service.OrderService
}

func NewOrdersHTTPHandler(srv service.OrderService) OrdersHTTPHandler {
	rpcHandler := &ordersHTTPHandler{
		service: srv,
	}

	return rpcHandler
}

func (h *ordersHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest

	if err := helpers.ParseJSON(r, &req); err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.CreateOrder(r.Context(), entity.NewOrder(1, 1, 1, 10))
	if err != nil {
		errors.WriteJSONError(w, http.StatusConflict, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
