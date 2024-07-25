package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/qwark97/another_ai_project/alog"
	"github.com/qwark97/another_ai_project/server/model"
)

type Actioner interface {
	Act(ctx context.Context, request model.Request) model.Response
}

type Server struct {
	router *mux.Router
	ai     Actioner
	log    alog.Logger
}

func NewServer(router *mux.Router, ai Actioner, log alog.Logger) Server {
	return Server{
		router: router,
		ai:     ai,
		log:    log,
	}
}

func (s Server) RegisterRoutes() error {
	s.router.HandleFunc("/api/v1/interaction", s.interaction).Methods("POST")
	s.router.HandleFunc("/api/v1/chat", s.chat)

	return nil
}

func (s Server) interaction(w http.ResponseWriter, r *http.Request) {
	s.log.Info("%s %s", r.Method, r.URL)
	now := time.Now()
	defer func(n time.Time) {
		s.log.Debug("%s %s took %s", r.Method, r.URL, time.Since(n).Round(time.Millisecond))
	}(now)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	var data model.Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := s.ai.Act(ctx, data)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Server) chat(w http.ResponseWriter, r *http.Request) {
	s.log.Info("%s (ws) %s", r.Method, r.URL)
	now := time.Now()
	defer func(n time.Time) {
		s.log.Debug("%s %s took %s", r.Method, r.URL, time.Since(n).Round(time.Millisecond))
	}(now)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	defer conn.Close()

	ctx := r.Context()

	for {
		var request model.Request
		err := conn.ReadJSON(&request)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Info("closing connection with %s", conn.RemoteAddr())
			} else {
				s.log.Error("failed to read the message: %s", err.Error())
			}
			return
		}
		func() {
			ctx, close := context.WithTimeout(ctx, 5*time.Minute)
			defer close()

			response := s.ai.Act(ctx, request)

			err = conn.WriteJSON(response)
			if err != nil {
				s.log.Error(err.Error())
				return
			}
		}()
	}
}
