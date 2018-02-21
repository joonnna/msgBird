package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/joonnna/msgbird"
)

func main() {
	var key string

	args := flag.NewFlagSet("args", flag.ExitOnError)
	args.StringVar(&key, "key", "0kMmexo2Q4gnQa2k4pZ2ZJxzO", "Access key to use when contacting the messagebird api")
	args.Parse(os.Args[1:])

	p, err := bird.NewProxy(key)
	if err != nil {
		panic(err)
	}

	p.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	p.Stop()
}
