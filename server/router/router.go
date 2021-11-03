package router

import (
	"ghotos/server/app"
	"ghotos/server/requestlog"

	"ghotos/server/app/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func New(a *app.App) *chi.Mux {
	l := a.Logger()
	r := chi.NewRouter()

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

	//r.Method("GET", "/", requestlog.NewHandler(a.HandleIndex, l))
	r.Method("GET", "/test", requestlog.NewHandler(a.HandleIndex2, l))
	//r.Get("/healthz/liveness", app.HandleLive)
	//r.Method("GET", "/healthz/readiness", requestlog.NewHandler(a.HandleReady, l))
	// Routes for healthz
	r.Get("/healthz", app.HandleHealth)

	// Routes for books
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)
		r.Method("POST", "/auth/login", requestlog.NewHandler(a.HandleAuthLogin, l))
		r.Method("POST", "/auth/refresh", requestlog.NewHandler(a.HandleAuthRefresh, l))

		// Routes for books
		r.Group(func(r chi.Router) {
			r.Method("GET", "/books", requestlog.NewHandler(a.HandleListBooks, l))
			r.Method("POST", "/books", requestlog.NewHandler(a.HandleCreateBook, l))
			r.Method("GET", "/books/{id}", requestlog.NewHandler(a.HandleReadBook, l))
			r.Method("PUT", "/books/{id}", requestlog.NewHandler(a.HandleUpdateBook, l))
			r.Method("DELETE", "/books/{id}", requestlog.NewHandler(a.HandleDeleteBook, l))
		})

		r.Method("GET", "/f/{src}", requestlog.NewHandler(a.HandleShowFile, l))

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(a))
			r.Method("GET", "/gallery", requestlog.NewHandler(a.HandleGallery, l))
			r.Method("GET", "/gallery/{day}", requestlog.NewHandler(a.HandleListGalleryDayFile, l))
			r.Method("POST", "/file", requestlog.NewHandler(a.HandleCreateFile, l))
			r.Method("DELETE", "/file/{uid}", requestlog.NewHandler(a.HandleDeleteFile, l))
			r.Method("GET", "/file/src/{uid}", requestlog.NewHandler(a.HandleReadFile, l))
		})

		r.Method("GET", "/auth/logout", requestlog.NewHandler(a.HandleAuthLogout, l))

	})

	r.Handle("/*", app.HandleClient())

	return r
}
