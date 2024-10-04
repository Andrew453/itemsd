package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"prjs/itemsd/app"
	"syscall"
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

func Manage() (err error) {

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
			return nil
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			allDone()
			if killSignal == os.Interrupt {
				return nil
			}
			return nil
		}
	}
}

func main() {
	err := Manage()
	if err != nil {
		errlog.Println("\nError: ", err)
		os.Exit(1)
	}
}
