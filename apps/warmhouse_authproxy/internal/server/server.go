package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/warmhouse/warmhouse_authproxy/internal/config"
	"github.com/warmhouse/warmhouse_authproxy/internal/models"
	"github.com/warmhouse/warmhouse_authproxy/internal/services/auth"
)

type Server struct {
	authService *auth.Service
	server      *http.Server
	client      *http.Client
	nextHost    string
}

func NewServer(nextHost string, conf *config.Config, authService *auth.Service) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.App.Port),
		Handler: mux,
	}

	s := &Server{server: server, nextHost: nextHost, client: http.DefaultClient, authService: authService}

	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/register", s.handleRegister)
	mux.HandleFunc("/", s.handleDefault)

	return s
}

func (s *Server) Start() error {
	log.Println("Starting server on", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleDefault(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling default request")

	userID, err := s.validateToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	url := fmt.Sprintf("http://%s/api/v1%s", s.nextHost, r.URL.Path)
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	log.Println("Proxying request to:", url)

	newRequest, err := http.NewRequestWithContext(r.Context(), r.Method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}

	log.Println("Setting X-User-Id to:", userID)
	newRequest.Header.Set("X-User-Id", userID)

	resp, err := s.client.Do(newRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	url := fmt.Sprintf("http://%s/api/v1%s", s.nextHost, r.URL.Path)
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	newRequest, err := http.NewRequestWithContext(r.Context(), r.Method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}

	resp, err := s.client.Do(newRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) validateToken(w http.ResponseWriter, r *http.Request) (string, error) {
	ctx := r.Context()

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return "", fmt.Errorf("token is required")
	}

	userID, err := s.authService.ValidateToken(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return "", fmt.Errorf("invalid token")
	}

	return userID, nil
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	url := fmt.Sprintf("http://%s/api/v1%s", s.nextHost, r.URL.Path)
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	log.Println("Proxying login request to:", url)

	newRequest, err := http.NewRequestWithContext(ctx, r.Method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println("Error creating new request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}

	resp, err := s.client.Do(newRequest)
	if err != nil {
		log.Println("Error doing request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("Response body:", string(body))

	var response models.UserAPIResponse

	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Error unmarshalling response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if response.User.ID == uuid.Nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("invalid credentials"))
		return
	}

	token, expiresAt, err := s.authService.GenerateToken(ctx, response.User)
	if err != nil {
		log.Println("Error generating token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
