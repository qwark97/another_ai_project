package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vargspjut/wlog"

	"github.com/gorilla/mux"
	"github.com/qwark97/assistant/server/agent"
	"github.com/qwark97/assistant/server/handlers"
	"github.com/qwark97/assistant/server/storage/data"
	"github.com/qwark97/assistant/server/storage/embedding"
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
	var a agent.Agent
	{
		d := data.New()

		e := embedding.New(log)
		if err := e.Connect(env["QDRANT_ADDR"]); err != nil {
			log.Fatal(err)
		}
		defer e.Disconnect()
		log.Info("connected to embedding engine")

		llmsGroup := agent.NewLLMsGroup(env, log)

		a = agent.New(d, e, llmsGroup, log)
	}

	router := mux.NewRouter()
	server := handlers.NewServer(router, a, log)
	if err := server.RegisterRoutes(); err != nil {
		return err
	}
	log.Info("registered routes")

	addr := fmt.Sprintf("%s:%s", conf.host, conf.port)
	log.Infof("starts listening at %s", addr)
	return http.ListenAndServe(addr, router)
}
