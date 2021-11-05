package router

import (
	"ghotos/server/app"

	"ghotos/server/app/middleware"

	"github.com/go-chi/chi"

	"github.com/go-chi/cors"
)

func New(a *app.App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger("", nil))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", a.HandleIndex)
	r.Get("/healthz", app.HandleHealth)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", a.HandleAuthLogin)
		r.Post("/auth/refresh", a.HandleAuthRefresh)

		r.Get("/f/{src}", a.HandleShowFile)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(a))
			r.Get("/auth/logout", a.HandleAuthLogout)
			r.Get("/gallery", a.HandleGallery)
			r.Get("/gallery/{day}", a.HandleListGalleryDayFile)
			r.Post("/file", a.HandleCreateFile)
			r.Delete("/file/{uid}", a.HandleDeleteFile)
			r.Get("/file/src/{uid}", a.HandleReadFile)
		})

	})
	/*
		//r.Method("GET", "/", a.HandleIndex)
		r.Method("GET", "/test", a.HandleIndex2)
		//r.Get("/healthz/liveness", app.HandleLive)
		//r.Method("GET", "/healthz/readiness", a.HandleReady)
		// Routes for healthz
		r.Get("/healthz", app.HandleHealth)

		// Routes for books
		r.Route("/api/v1", func(r chi.Router) {
			r.Use(middleware.ContentTypeJson)
			r.Method("POST", "/auth/login", a.HandleAuthLogin)
			r.Method("POST", "/auth/refresh", a.HandleAuthRefresh)

			// Routes for books
			r.Group(func(r chi.Router) {
				r.Method("GET", "/books", a.HandleListBooks)
				r.Method("POST", "/books", a.HandleCreateBook)
				r.Method("GET", "/books/{id}", a.HandleReadBook)
				r.Method("PUT", "/books/{id}", a.HandleUpdateBook)
				r.Method("DELETE", "/books/{id}", a.HandleDeleteBook)
			})

			r.Method("GET", "/f/{src}", a.HandleShowFile)

			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(a))
				r.Method("GET", "/gallery", a.HandleGallery)
				r.Method("GET", "/gallery/{day}", a.HandleListGalleryDayFile)
				r.Method("POST", "/file", a.HandleCreateFile)
				r.Method("DELETE", "/file/{uid}", a.HandleDeleteFile)
				r.Method("GET", "/file/src/{uid}", a.HandleReadFile)
			})

			r.Method("GET", "/auth/logout", a.HandleAuthLogout)

		})
	*/

	r.Handle("/*", app.HandleClient())

	return r
}
