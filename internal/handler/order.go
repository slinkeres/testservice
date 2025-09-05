package handler

import (
	"encoding/json"
	"net/http"
	"order-service/internal/cache"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	cache *cache.Cache
}

func NewOrderHandler(cache *cache.Cache) *OrderHandler {
	return &OrderHandler{cache: cache}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]

	if orderUID == "" {
		http.Error(w, "Требуется ID заказа", http.StatusNotFound)
		return
	}

	order, exists := h.cache.Get(orderUID)

	if !exists {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
	}
}

func (h *OrderHandler)  RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/order/{id}", h.GetOrder).Methods("Get")
}