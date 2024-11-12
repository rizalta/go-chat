package server

import (
	"context"
	"go-chat/cmd/web"
	"go-chat/cmd/web/components"
	"go-chat/cmd/web/pages"
	"go-chat/internal/handler"
	"go-chat/internal/middlewares"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewares.AuthMiddleware)

	userHandler := handler.NewUserHandler(s.userRepo)
	wsHandler := handler.NewWSHandler(s.hub)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/", s.serveIndex)
	r.Get("/login", userHandler.ServeLogin)
	r.Get("/signup", userHandler.ServeSignup)

	r.Post("/login", userHandler.Login)
	r.Post("/signup", userHandler.Signup)
	r.Get("/signout", userHandler.Signout)

	r.Get("/ws", wsHandler.HandleWS(ctx))

	r.Get("/messages/{page}", s.loadMessages(ctx))

	return r
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	err := pages.Index(userID).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func (s *Server) loadMessages(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		if r.Header.Get("HX-Request") != "true" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		page, _ := strconv.Atoi(r.PathValue("page"))
		messages, hasNext, _ := s.messageRepo.GetMessages(ctx, int64(page))
		components.LoadMessages(messages, userID, page, hasNext).Render(ctx, w)
	}
}
