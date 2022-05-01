package main

import (
	"HELLOWORD/handlers"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	serverAddress string = ":3000"
)

func routes() (r *chi.Mux) {
	r = chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestID,
			middleware.RealIP,
			middleware.Logger)
		r.Get("/health", handlers.Health)
		r.Get("/metrics", promhttp.Handler().ServeHTTP)
		r.Get("/load", handlers.Load)
		r.Get("/apis/v1/deployments", handlers.Deployments)
		r.Get("/apis/v1/deployment/{deploy:[0-9a-zA-Z-]+}", handlers.Pods)
		r.Get("/apis/v1/pods*", handlers.Podstatus)
		r.Get("/", handlers.Hello)
	})
	return r
}

func main() {
	var (
		r   *chi.Mux
		s   *http.Server
		err error
	)
	r = routes()
	s = &http.Server{
		Addr:         serverAddress,
		Handler:      r,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err = s.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	fmt.Println("Server started at" + serverAddress)

	// Gracefully shutdown the server if we recieve a signal
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println(err.Error())
	}
}
