module github.com/qwark97/another_ai_project/server

go 1.21.3

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/qwark97/another_ai_project/alog v0.0.0-00010101000000-000000000000
	github.com/qwark97/another_ai_project/llms v0.0.0-00010101000000-000000000000
)

replace github.com/qwark97/another_ai_project/llms => ../llms

replace github.com/qwark97/another_ai_project/alog => ../alog
