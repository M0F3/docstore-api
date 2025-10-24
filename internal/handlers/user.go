package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/M0F3/docstore-api/internal/models"
	"github.com/M0F3/docstore-api/internal/services"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserPayload
	w.Header().Set("Content-Type", "Application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(err) 
		return 
	}
	registered_user, err := h.Service.Register(r.Context(), req)
	log.Println("Registered")

	if err != nil {
		json.NewEncoder(w).Encode(err) 
		return 
	}
	json.NewEncoder(w).Encode(registered_user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginPayload
	w.Header().Set("Content-Type", "Application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(err) 
		return 
	}
	result, err := h.Service.Login(r.Context(), req)

	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(err) 
		return 
	}

	w.Header().Add("Authtoken", result.Token)
	json.NewEncoder(w).Encode(result.User)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	l, err := h.Service.ListUsers(r.Context())
	if err != nil {
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(l)
}