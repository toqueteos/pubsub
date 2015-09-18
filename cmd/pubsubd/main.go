package pubsub

import (
	"flag"
	"os"
	"os/signal"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/ikanor/pubsub/server"
)

var addr string

func main() {
	flag.Parse()

	flag.StringVar(&addr, "addr", "0.0.0.0", "")

	pubsubd, err := server.New(addr)
	if err != nil {
		log15.Info("server.New", "error", err)
		os.Exit(1)
	}

	go pubsubd.Start()
	log15.Info("Listening", "on", addr)

	go pubsubd.ForEachError(func(err error) {
		if err != nil {
			log15.Info("pubsubd", "error", err)
		}
	})

	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt)
	<-ctrlc
	pubsubd.Stop()
}
