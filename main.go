package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strings"
	"whiteListApi/routes/country"
)

// setup base routing and default router settings
// in production some of this may come from configuration, for instance the timeout
// see https://github.com/go-chi/chi#middlewares for info on middleware
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.Heartbeat("/v1/api/ping"),
		middleware.Timeout(15),
	)

	// just for fun throw some cors headers, also using middleware
	corsSettings := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(corsSettings.Handler)

	// Simple versioning scheme
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/country", country.Routes())
		// r.Mount("/api/foo", foo.Routes())  //can easily add new branches to route, idiomatic, extensible
	})

	return router
}

func main() {
	router := Routes()
	port := ":8080" // NOTE port would come from configuration in production system

	log.Printf("Starting server on port %s", port)
	log.Printf("%5v: %s \n", "GET", "/v1/api/ping") // chi.Walk doesn't walk the heartbeat route

	// useful to walk the routes and print them, especially during development
	wFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		log.Printf("%5v: %s \n", method, strings.Replace(route, "*/", "", -1))
		return nil
	}
	if err := chi.Walk(router, wFunc); err != nil {
		log.Panicf("Error walking routes: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(port, router))
}
