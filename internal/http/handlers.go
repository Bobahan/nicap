package app_http

import (
	"artofdev/data"
	"artofdev/entity"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handlers struct {
	log *log.Logger
}

func New(log *log.Logger) *Handlers {
	return &Handlers{log}
}

func (h *Handlers) GetProducts(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetProducts"
	if uuid, ok := r.Context().Value(RequestIdKey).(string); ok {
		log.Printf("Обрабатывается %s запрос с UUID: %s\n", op, uuid)
	}

	lp := data.GetProducts()

	if err := lp.ToJson(w); err != nil {
		http.Error(w, "Failed to parse json", http.StatusInternalServerError)
		return
	}

	log.Println("успешно отравили данные")
}

// AddProduct godoc
// @Tags AddProduct v1
// @Summary AddProduct
// @Description Ручка принимает json, приводит к типу json и добавляет 1 products в products
// @Accept json
// @Produce json
// @Param request body entity.ApiRequest true "Данные продукта"
// @Success 200 {object} entity.ApiResponse
// @Router /api/v1/products [post]
func (h *Handlers) AddProduct(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.AddProduct"
	if uuid, ok := r.Context().Value(RequestIdKey).(string); ok {
		log.Printf("Обрабатывается %s запрос с UUID: %s\n", op, uuid)
	}

	product := r.Context().Value(RequestProductKey).(*data.Product)

	products := data.AddProduct(product)

	response := entity.ApiResponse{
		Success: true,
		Message: "Продукт был успешно добавлен",
		Data:    products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetProducts"
	if uuid, ok := r.Context().Value(RequestIdKey).(string); ok {
		log.Printf("Обрабатывается %s запрос с UUID: %s\n", op, uuid)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	product := r.Context().Value(RequestProductKey).(*data.Product)

	err := data.UpdateProduct(product, id)
	if errors.Is(err, data.ErrProductNotFound) {
		http.Error(w, "Failed to find product by id", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Failed to find product by id", http.StatusInternalServerError)
		return
	}
}
