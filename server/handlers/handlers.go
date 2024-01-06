package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qwark97/assistant/server/controller"
	"github.com/qwark97/assistant/server/model"
)

type Server struct {
	router *mux.Router
	cont   controller.Controller
	log    *slog.Logger
}

func NewServer(router *mux.Router, log *slog.Logger) Server {
	cont := controller.New(log)
	return Server{
		router: router,
		cont:   cont,
		log:    log,
	}
}

func (s Server) RegisterRoutes() error {
	s.router.HandleFunc("/api/v1/interaction", s.interaction).Methods("POST")

	return nil
}

func (s Server) interaction(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("received")
	w.Header().Add("Content-Type", "application/json")

	var data model.ActionRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.log.Error(fmt.Sprintf("decode: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	llm := defaultLLM()

	category, typ := s.cont.RecogniseInteraction(data.Instruction, llm)
	_ = category
	_ = typ

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.log.Error(fmt.Sprintf("encode: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
