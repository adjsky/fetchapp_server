package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adjsky/fetchapp_server/internal/application"
)

func main() {
	log.SetFlags(0)
	app := application.New()
	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeChan
		app.Close()
		os.Exit(0)
	}()
	app.Start()
}
