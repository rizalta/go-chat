package server

import (
	"go-chat/cmd/web"
	"go-chat/cmd/web/pages"
	"go-chat/internal/handler"
	"go-chat/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewares.AuthMiddleware)

	userHandler := handler.NewUserHandler(s.db)
	wsHandler := handler.NewWSHandler(s.hub)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("userID").(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		pages.Index().Render(r.Context(), w)
	})
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		pages.Login().Render(r.Context(), w)
	})
	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		pages.Signup().Render(r.Context(), w)
	})

	r.Post("/login", userHandler.Login)
	r.Post("/signup", userHandler.Signup)
	r.Get("/signout", userHandler.Signout)

	r.Get("/ws", wsHandler.HandleWS)

	return r
}
