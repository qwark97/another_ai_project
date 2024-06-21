package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vargspjut/wlog"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/qwark97/another_ai_project/server/model"
)

type Actioner interface {
	Act(ctx context.Context, request model.Request) model.Response
}

type Server struct {
	router *mux.Router
	ai     Actioner
	log    wlog.Logger
}

func NewServer(router *mux.Router, ai Actioner, log wlog.Logger) Server {
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	var data model.Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := s.ai.Act(ctx, data)
	s.log.Debug("response:", response)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Server) chat(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s.log.Infof("received %s request", r.Method)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("upgrader: %v", err)
		return
	}
	defer conn.Close()

	ctx := r.Context()

	for {
		var request model.Request
		err := conn.ReadJSON(&request)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Info("closing connection")
			} else {
				s.log.Errorf("failed to read the message: %s", err.Error())
			}
			return
		}
		func() {
			ctx, close := context.WithTimeout(ctx, 5*time.Minute)
			defer close()

			response := s.ai.Act(ctx, request)

			s.log.Debug("response:", response)
			err = conn.WriteJSON(response)
			if err != nil {
				s.log.Errorf("failed to write the response message: %s", err.Error())
				return
			}
		}()
	}
}
