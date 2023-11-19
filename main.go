package main

import (
	"github.com/dwnGnL/ddos-pow/cmd"
	"github.com/dwnGnL/ddos-pow/config"
	"log"
	"os"

	"github.com/dwnGnL/ddos-pow/lib/goerrors"

	"github.com/sirupsen/logrus"
)

var Version = "v0.0.1"

type command string

const (
	CLIENT command = "client" // command to start the client service
	SERVER command = "server" // command to start the server service
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("service must be runned with one parameter!")
		return
	}

	cfg := config.FromFile(os.Getenv("CONFIG_FILE"))
	intLogger(cfg.LogLevel)
	service := command(args[1])
	
	switch service {
	case CLIENT:
		err := cmd.StartClient(cfg)
		if err != nil {
			log.Fatalf("app run: %s", err)
		}
	case SERVER:
		err := cmd.StartServer(cfg)
		if err != nil {
			log.Fatalf("app run: %s", err)
		}
	}
}

func intLogger(logLevel string) {
	var formatter logrus.Formatter = new(logrus.JSONFormatter)
	if os.Getenv("LOG_FORMAT") == "text" {
		formatter = new(logrus.TextFormatter)
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	err = goerrors.Setup(formatter, level)
	if err != nil {
		panic(err)
	}
}
