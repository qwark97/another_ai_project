module github.com/qwark97/assistant/server

go 1.21.3

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/vargspjut/wlog v1.0.11
)

replace github.com/qwark97/assistant/llms => ../llms
