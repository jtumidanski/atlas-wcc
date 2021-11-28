package main

import (
   "atlas-wcc/character/instruction"
   "atlas-wcc/configuration"
   "atlas-wcc/kafka/consumers"
   "atlas-wcc/kafka/producers"
   "atlas-wcc/logger"
   "atlas-wcc/rest"
   "atlas-wcc/session"
   "atlas-wcc/socket"
   "atlas-wcc/tracing"
   "context"
   "io"
   "os"
   "os/signal"
   "strconv"
   "sync"
   "syscall"
)

const serviceName = "atlas-wcc"

func main() {
   l := logger.CreateLogger(serviceName)
   l.Infoln("Starting main service.")

   wg := &sync.WaitGroup{}
   ctx, cancel := context.WithCancel(context.Background())

   tc, err := tracing.InitTracer(l)(serviceName)
   if err != nil {
      l.WithError(err).Fatal("Unable to initialize tracer.")
   }
   defer func(tc io.Closer) {
      err := tc.Close()
      if err != nil {
         l.WithError(err).Errorf("Unable to close tracer.")
      }
   }(tc)

   _, err = configuration.Get()
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

   rest.CreateService(l, ctx, wg, "/ms/csrv/worlds/{worldId}/channels/{channelId}", session.InitResource, instruction.InitResource)

   sl, span := tracing.StartSpan(l, "startup")
   producers.StartChannelServer(sl, span)(byte(wid), byte(cid), ha, uint32(port))
   span.Finish()

   // trap sigterm or interrupt and gracefully shutdown the server
   c := make(chan os.Signal, 1)
   signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

   // Block until a signal is received.
   sig := <-c
   l.Infof("Initiating shutdown with signal %s.", sig)
   cancel()
   wg.Wait()

   sl, span = tracing.StartSpan(l, "shutdown")
   producers.ShutdownChannelServer(sl, span)(byte(wid), byte(cid), ha, uint32(port))
   session.DestroyAll(sl, span, session.Registry())
   span.Finish()

   l.Infoln("Service shutdown.")
}
