package pubsub

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/ikanor/pubsub/client"
	"github.com/ikanor/pubsub/server"
)

var addr string

func main() {
	flag.Parse()

	flag.StringVar(&addr, "addr", "0.0.0.0", "")

	pubsubd, err := server.New(addr)
	if err != nil {
		log15.Error("server.New", "error", err)
		os.Exit(1)
	}

	go pubsubd.Start()
	log15.Info("Listening", "on", addr)

	go pubsubd.ForEachError(func(err error) {
		if err != nil {
			log15.Error("pubsubd", "error", err)
		}
	})

	var clients []client.Client
	for i := 0; i < 4; i++ {
		c, err := client.New("127.0.0.1")
		if err != nil {
			log15.Error("client.New", "error", err)
			os.Exit(1)
		}
		c.Subscribe("test", func(c client.Client, payload []byte) {
			fmt.Printf("client-%x %q\n", c.ID, payload)
		})
	}

	pubsubd.Publish("test", "foo")
	pubsubd.Publish("test", "bar")
	pubsubd.Publish("test", "qux")

	fmt.Println("You should see 12 messages, press Ctrl+C to close this demo")

	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt)
	<-ctrlc
	pubsubd.Stop()
}
