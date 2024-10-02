package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"prjs/itemsd/app"
	"syscall"

	"github.com/takama/daemon"
)

const (
	name        = "items_daemon"
	description = "items service"
)

var (
	stdlog, errlog *log.Logger
	doneContext    context.Context
	doneFunc       context.CancelFunc
)

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (usage string, err error) {
	usage = fmt.Sprintf("Usage: %s install | remove | start | stop | status", name)

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	doneContext, doneFunc = context.WithCancel(context.Background())
	a, err := app.New(doneContext)
	if err != nil {
		stdlog.Println(err)
		errlog.Println(err)
		return
	}
	errorsChan := a.Run()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	allDone := func() {
		err = a.Stop()
		if err != nil {
			stdlog.Println(err)
			errlog.Println(err)
		}
	}

	for {
		select {
		case err = <-errorsChan:
			stdlog.Println(err)
			errlog.Println(err)
			allDone()
			return "Daemon has died by himself", nil
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			allDone()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

func main() {
	srv, err := daemon.New(name, description, daemon.SystemDaemon)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}

	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
}
