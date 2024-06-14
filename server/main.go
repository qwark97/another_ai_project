package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vargspjut/wlog"

	"github.com/gorilla/mux"
	"github.com/qwark97/another_ai_project/server/ai"
	"github.com/qwark97/another_ai_project/server/handlers"
)

func main() {
	startingConfiguration := parseFlags()
	envConfiguration := loadEnvFile(startingConfiguration.envFilePath)
	log := newLogger()

	if err := run(startingConfiguration, envConfiguration, log); err != nil {
		panic(err)
	}
	log.Info("stopping service")

}

func newLogger() wlog.Logger {
	return wlog.New(os.Stdout, wlog.Dbg, false)
}

func run(conf *flagsConf, env environmentVars, log wlog.Logger) error {
	router := mux.NewRouter()

	ai := ai.New()
	server := handlers.NewServer(router, ai, log)
	if err := server.RegisterRoutes(); err != nil {
		return err
	}
	log.Info("registered routes")

	addr := fmt.Sprintf("%s:%s", conf.host, conf.port)
	log.Infof("starts listening at %s", addr)
	return http.ListenAndServe(addr, router)
}
