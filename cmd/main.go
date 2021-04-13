package main

import (
	"os"
	"os/signal"
	"server/internal/application"
	"syscall"
)

func main() {
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
