module github.com/qwark97/assistant/server

go 1.21.3

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/qdrant/go-client v1.9.0
	github.com/qwark97/assistant/llms v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
	github.com/vargspjut/wlog v1.0.11
	golang.org/x/sync v0.6.0
	google.golang.org/grpc v1.64.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/qwark97/assistant/llms => ../llms
