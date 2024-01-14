package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vargspjut/wlog"

	"github.com/gorilla/mux"
	"github.com/qwark97/assistant/llms/openai"
	"github.com/qwark97/assistant/server/controller"
	"github.com/qwark97/assistant/server/model"
)

type Server struct {
	router *mux.Router
	cont   controller.Controller
	log    wlog.Logger
	env    map[string]string
}

func NewServer(router *mux.Router, env map[string]string, log wlog.Logger) Server {
	cont := controller.New(log)
	return Server{
		router: router,
		cont:   cont,
		log:    log,
		env:    env,
	}
}

func (s Server) RegisterRoutes() error {
	s.router.HandleFunc("/api/v1/interaction", s.interaction).Methods("POST")

	return nil
}

func (s Server) interaction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*5)
	defer cancel()

	s.log.Debug("received")
	w.Header().Add("Content-Type", "application/json")

	var data model.ActionRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.log.Error(fmt.Sprintf("decode: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	llm := openai.New(s.env["OPENAI_KEY"], s.log)

	response := s.cont.Interact(ctx, data.Instruction, llm)

	res := map[string]string{
		"response": response,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error(fmt.Sprintf("encode: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
