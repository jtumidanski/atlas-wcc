package main

import (
	"atlas-wcc/kafka/consumers"
	"atlas-wcc/kafka/producers"
	"atlas-wcc/logger"
	"atlas-wcc/registries"
	"atlas-wcc/rest"
	"atlas-wcc/services"
	"atlas-wcc/socket"
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	_, err := registries.GetConfiguration()
	if err != nil {
		l.WithError(err).Fatalf("Unable to successfully load configuration.")
	}

	wid, err := strconv.ParseUint(os.Getenv("WORLD_ID"), 10, 8)
	if err != nil {
		l.WithError(err).Fatalf("Unable to read world identifier from environment.")
		return
	}
	cid, err := strconv.ParseUint(os.Getenv("CHANNEL_ID"), 10, 8)
	if err != nil {
		l.WithError(err).Fatalf("Unable to read channel identifier from environment.")
		return
	}
	ha := os.Getenv("HOST_ADDRESS")
	port, err := strconv.ParseUint(os.Getenv("CHANNEL_PORT"), 10, 32)
	if err != nil {
		l.WithError(err).Fatalf("Unable to read port from environment.")
		return
	}

	consumers.CreateEventConsumers(l, byte(wid), byte(cid), ctx, wg)

	socket.CreateSocketService(l, ctx, wg)(byte(wid), byte(cid), int(port))

	rest.CreateRestService(l, ctx, wg)

	producers.StartChannelServer(l)(byte(wid), byte(cid), ha, uint32(port))

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()

	producers.ShutdownChannelServer(l)(byte(wid), byte(cid), ha, uint32(port))
	services.DestroyAll(l, registries.GetSessionRegistry())

	l.Infoln("Service shutdown.")
}
