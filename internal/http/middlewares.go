package app_http

import (
	"artofdev/data"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	RequestIdKey      contextKey = "requestID"
	RequestProductKey contextKey = "product"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func SetRequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New().String()

		w.Header().Set("X-REQUEST-ID", uuid)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIdKey, uuid)))
	})
}

// Чтобы не повторять каждый раз чтение данных из req.Body и приведение в json
// Мы вынесли эту логику в Middleware
func ValidateReqBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := &data.Product{}

		// читаем из запроса и загоняем в product
		if err := product.FromJson(r.Body); err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// валидируем объект product
		if err := product.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), RequestProductKey, product)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
