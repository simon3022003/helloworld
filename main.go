package main

import (
	"technical_challenge/handlers"
	"os"
	"os/signal"
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
	"context"
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
