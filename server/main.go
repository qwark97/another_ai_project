package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/qwark97/assistant/server/handlers"
)

func main() {
	startingConfiguration := parseFlags()
	log := newLogger()

	if err := run(startingConfiguration, log); err != nil {
		panic(err)
	}
	log.Info("stopping service")

}

func newLogger() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})
	return slog.New(jsonHandler)
}

func run(conf *flagsConf, log *slog.Logger) error {
	router := mux.NewRouter()

	server := handlers.NewServer(router, log)
	if err := server.RegisterRoutes(); err != nil {
		return err
	}
	log.Info("registered routes")

	addr := fmt.Sprintf("%s:%s", conf.host, conf.port)
	log.Info("starts listening", "addr", addr)
	return http.ListenAndServe(addr, router)
}
