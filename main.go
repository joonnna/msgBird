package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/joonnna/msgBird"
)

func main() {
	var key string

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&key, "key", "JhVoK80WnKAEk8je4e8RgUykG", "Access key to use when contacting the messagebird api")
	args.Parse(os.Args[1:])

	p, err := bird.NewProxy(key)
	if err != nil {
		panic(err)
	}

	p.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	<-ch

	p.Stop()
}
