package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/a1uka/rzhaka_tournaments/internal/app"
)

func main() {

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	application.Shutdown()
}
