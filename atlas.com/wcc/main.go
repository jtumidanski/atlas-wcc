package main

import (
	"atlas-wcc/kafka/consumers"
	"atlas-wcc/kafka/producers"
	"atlas-wcc/logger"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/registries"
	"atlas-wcc/rest"
	"atlas-wcc/services"
	"atlas-wcc/socket/request"
	"atlas-wcc/socket/request/handler"
	"context"
	"github.com/jtumidanski/atlas-socket"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	l := logger.CreateLogger()

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

	consumers.CreateEventConsumers(l, byte(wid), byte(cid))

	createSocketService(l, wid, cid, err, port)

	rest.CreateRestService(l)

	producers.ChannelServer(l, context.Background()).Start(byte(wid), byte(cid), ha, uint32(port))

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infoln("Shutting down via signal:", sig)
	producers.ChannelServer(l, context.Background()).Shutdown(byte(wid), byte(cid), ha, uint32(port))

	processors.ForEachSession(disconnect())
}

func disconnect() processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Disconnect()
	}
}

func createSocketService(l *logrus.Logger, wid uint64, cid uint64, err error, port uint64) {
	lss := services.NewMapleSessionService(l, byte(wid), byte(cid))

	w := l.Writer()
	defer w.Close()

	ss, err := socket.NewServer(log.New(w, "", 0), lss, socket.IpAddress("0.0.0.0"), socket.Port(int(port)))
	if err != nil {
		l.Fatal(err.Error())
	}

	registerSocketRequestHandlers(ss, l)
	go ss.Run()
}

func registerSocketRequestHandlers(ss *socket.Server, l logrus.FieldLogger) {
	hr := func(op uint16, validator request.SessionStateValidator, handler request.SessionRequestHandler) {
		ss.RegisterHandler(op, request.AdaptHandler(l, validator, handler))
	}
	hr(handler.OpCodePong, request.NoOpValidator(), handler.PongHandler())
	hr(handler.OpCharacterLoggedIn, request.NoOpValidator(), handler.CharacterLoggedInHandler())
	hr(handler.OpChangeMapSpecial, request.LoggedInValidator(), handler.ChangeMapSpecialHandler())
	hr(handler.OpMoveCharacter, request.LoggedInValidator(), handler.MoveCharacterHandler())
	hr(handler.OpChangeMap, request.LoggedInValidator(), handler.ChangeMapHandler())
	hr(handler.OpMoveLife, request.LoggedInValidator(), handler.MoveLifeHandler())
	hr(handler.OpGeneralChat, request.LoggedInValidator(), handler.GeneralChatHandler())
	hr(handler.OpChangeChannel, request.LoggedInValidator(), handler.ChangeChannelHandler())
	hr(handler.OpCharacterExpression, request.LoggedInValidator(), handler.CharacterExpressionHandler())
	hr(handler.OpCharacterCloseRangeAttack, request.LoggedInValidator(), handler.CharacterCloseRangeAttackHandler())
	hr(handler.OpCharacterDistributeAp, request.LoggedInValidator(), handler.DistributeApHandler())
	hr(handler.OpCharacterDistributeSp, request.LoggedInValidator(), handler.DistributeSpHandler())
	hr(handler.OpCharacterHealOverTime, request.LoggedInValidator(), handler.HealOverTimeHandler())
	hr(handler.OpCharacterItemPickUp, request.LoggedInValidator(), handler.ItemPickUpHandler())
	hr(handler.OpNpcAction, request.LoggedInValidator(), handler.HandleNPCAction())
	hr(handler.OpNpcTalkMore, request.LoggedInValidator(), handler.HandleNPCTalkMoreRequest())
	hr(handler.OpNpcTalk, handler.CharacterAliveValidator(), handler.HandleNPCTalkRequest())
	hr(handler.OpCharacterDamage, request.LoggedInValidator(), handler.HandleCharacterDamageRequest())
}
