package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vargspjut/wlog"

	"github.com/gorilla/mux"
	"github.com/qwark97/assistant/server/agent"
	"github.com/qwark97/assistant/server/model"
)

type Server struct {
	router *mux.Router
	agent  agent.Agent
	log    wlog.Logger
}

func NewServer(router *mux.Router, agent agent.Agent, log wlog.Logger) Server {
	return Server{
		router: router,
		agent:  agent,
		log:    log,
	}
}

func (s Server) RegisterRoutes() error {
	s.router.HandleFunc("/api/v1/interaction", s.interaction).Methods("POST")
	s.router.HandleFunc("/api/v1/chat", s.chat)

	return nil
}

func (s Server) interaction(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("received")
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*5)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	var data model.InteractionRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.log.Error(fmt.Sprintf("decode: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if data.ConversationID == uuid.Nil {
		data.ConversationID = uuid.New()
	}

	response := s.agent.Interact(ctx, data)

	res := model.InteractionResponse{
		Answer:         <-response,
		ConversationID: data.ConversationID,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error(fmt.Sprintf("encode: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.log.Debug("responded")
}

func (s Server) chat(w http.ResponseWriter, r *http.Request) {
	// var upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// 	CheckOrigin: func(r *http.Request) bool {
	// 		return true
	// 	},
	// }
	// s.log.Infof("received %s request", r.Method)

	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	s.log.Error(fmt.Sprintf("upgrade: %s", err.Error()))
	// 	return
	// }
	// defer conn.Close()

	// llm := openai.New(s.env["OPENAI_KEY"], s.log)
	// ctx := r.Context()

	// for {
	// 	messageType, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
	// 			s.log.Info("closing connection")
	// 		} else {
	// 			s.log.Errorf("failed to read the message: %s", err.Error())
	// 		}
	// 		return
	// 	}
	// 	if messageType != websocket.TextMessage {
	// 		if err := conn.WriteMessage(messageType, []byte("Sorry, I can't handle this type of message")); err != nil {
	// 			s.log.Errorf("failed to write the message: %s", err.Error())
	// 			return
	// 		}
	// 	}

	// 	var data model.InteractionRequest
	// 	err = json.Unmarshal(p, &data)
	// 	if err != nil {
	// 		s.log.Errorf("failed to decode the message: %s", err.Error())
	// 		return
	// 	}

	// 	response := s.agent.Interact(ctx, data)

	// 	res := map[string]string{
	// 		"response":       response,
	// 		"converation_id": data.ConversationID.String(),
	// 	}
	// 	r, err := json.Marshal(res)
	// 	if err != nil {
	// 		s.log.Errorf("failed to encode the response message: %s", err.Error())
	// 		return
	// 	}

	// 	if err := conn.WriteMessage(messageType, r); err != nil {
	// 		s.log.Errorf("failed to write the response message: %s", err.Error())
	// 		return
	// 	}
	// }
}
