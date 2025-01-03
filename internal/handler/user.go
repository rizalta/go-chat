package handler

import (
	"context"
	"fmt"
	"go-chat/cmd/web/components"
	"go-chat/cmd/web/pages"
	"go-chat/internal/database"
	"go-chat/internal/utils"
	"net/http"
	"time"
)

type UserHandler struct {
	repo *database.UserRepo
}

func NewUserHandler(repo *database.UserRepo) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := utils.ValidateEmail(email); err != nil {
		sendErroNotification(r.Context(), w, err.Error())
		return
	}
	user, err := h.repo.GetUserByEmail(r.Context(), email)
	if err != nil {
		sendErroNotification(r.Context(), w, err.Error())
		return
	}

	if utils.CheckPasswordHash(user.Password, password) {
		sendErroNotification(r.Context(), w, "Invalid Password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		sendErroNotification(r.Context(), w, "Something went wrong")
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("HX-Redirect", "/")
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := utils.ValidateEmail(email); err != nil {
		sendErroNotification(r.Context(), w, err.Error())
		return
	}

	if err := utils.ValidatePassword(password); err != nil {
		sendErroNotification(r.Context(), w, err.Error())
		return
	}

	user, err := h.repo.AddUser(r.Context(), database.AddUserParams{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		sendErroNotification(r.Context(), w, err.Error())
		return
	}
	fmt.Println(user)
	w.Header().Set("HX-Redirect", "/login")
}

func (h *UserHandler) Signout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "session",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/login")
}

func sendErroNotification(ctx context.Context, w http.ResponseWriter, msg string) {
	components.ErrorNotification(msg).Render(ctx, w)
}

func (h *UserHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("userID").(string)
	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err := pages.Login().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func (h *UserHandler) ServeSignup(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("userID").(string)
	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err := pages.Signup().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}
