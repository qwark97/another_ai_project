package main

import "flag"

type flagsConf struct {
	host        string
	port        string
	envFilePath string
}

func parseFlags() *flagsConf {
	defer flag.Parse()
	var conf = new(flagsConf)

	flag.StringVar(&conf.host, "host", "127.0.0.1", "http server host value")
	flag.StringVar(&conf.port, "port", "8080", "http server port value")
	flag.StringVar(&conf.envFilePath, "env", "../.env", "path to .env file")
	return conf
}
