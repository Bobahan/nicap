package main

import (
	app_http "artofdev/internal/http"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "artofdev/docs"

	muxHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	log := log.New(os.Stdout, "product-api ", log.LstdFlags)

	handlers := app_http.New(log)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(muxHandler.CORS(muxHandler.AllowedOrigins([]string{"*"})))
	api.Use(app_http.LoggingMiddleware)
	api.Use(app_http.SetRequestIdMiddleware)
	api.Path("/products").Methods("GET").HandlerFunc(handlers.GetProducts)

	products := api.PathPrefix("/products").Subrouter()
	products.Use(app_http.ValidateReqBodyMiddleware)

	products.Path("").Methods("POST").HandlerFunc(handlers.AddProduct)
	products.Path("/{id}").Methods("POST").HandlerFunc(handlers.UpdateProduct)

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		IdleTimeout:  time.Second * 120,
		ReadTimeout:  time.Second * 4,
		WriteTimeout: time.Second * 4,
	}

	go func() {
		log.Println("Server starting on :8080")
		log.Println("Swagger UI: http://localhost:8080/swagger/index.html")
		log.Println("API Base URL: http://localhost:8080/api/v1")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("pid:", os.Getpid())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)
}
