package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vargspjut/wlog"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/qwark97/assistant/llms/openai"
	"github.com/qwark97/assistant/server/controller"
	"github.com/qwark97/assistant/server/model"
	"github.com/qwark97/assistant/server/storage/data"
)

type Server struct {
	router *mux.Router
	cont   controller.Controller
	log    wlog.Logger
	env    map[string]string
}

func NewServer(router *mux.Router, env map[string]string, log wlog.Logger) Server {
	cont := controller.New(data.New(), log)
	return Server{
		router: router,
		cont:   cont,
		log:    log,
		env:    env,
	}
}

func (s Server) RegisterRoutes() error {
	s.router.HandleFunc("/api/v1/interaction", s.interaction).Methods("POST")
	s.router.HandleFunc("/api/v1/chat", s.chat)

	return nil
}

func (s Server) interaction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*5)
	defer cancel()

	s.log.Debug("received")
	w.Header().Add("Content-Type", "application/json")

	var data model.InteractionRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.log.Error(fmt.Sprintf("decode: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	llm := openai.New(s.env["OPENAI_KEY"], s.log)

	response := s.cont.Interact(ctx, data, llm)

	res := map[string]string{
		"response": response,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error(fmt.Sprintf("encode: %s", err.Error()))
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
		s.log.Error(fmt.Sprintf("upgrade: %s", err.Error()))
		return
	}
	defer conn.Close()

	// llm := openai.New(s.env["OPENAI_KEY"], s.log)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			s.log.Error(err)
			return
		}

		s.log.Infof("recv: %s", p)
		s.log.Infof("type: %v", messageType)

		respnse := fmt.Sprintf("echo: %s", p)
		if err := conn.WriteMessage(messageType, []byte(respnse)); err != nil {
			s.log.Error(err)
			return
		}

		if string(p) == "exit" {
			break
		}
	}
	s.log.Info("done")
}
