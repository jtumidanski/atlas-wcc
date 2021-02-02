package main

import (
	"atlas-wcc/kafka/consumers"
	"atlas-wcc/kafka/producers"
	"atlas-wcc/registries"
	"atlas-wcc/rest"
	"atlas-wcc/services"
	"atlas-wcc/socket/request"
	"atlas-wcc/socket/request/handler"
	"context"
	"fmt"
	"github.com/jtumidanski/atlas-socket"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	l := log.New(os.Stdout, "wcc ", log.LstdFlags|log.Lmicroseconds)

	_, err := registries.GetConfiguration()
	if err != nil {
		l.Fatal("[ERROR] Unable to successfully load configuration.")
	}

	wid, err := strconv.ParseUint(os.Getenv("WORLD_ID"), 10, 8)
	if err != nil {
		l.Fatal("[ERROR] Unable to read world identifier from environment.")
		return
	}
	cid, err := strconv.ParseUint(os.Getenv("CHANNEL_ID"), 10, 8)
	if err != nil {
		l.Fatal("[ERROR] Unable to read channel identifier from environment.")
		return
	}
	ha := os.Getenv("HOST_ADDRESS")
	port, err := strconv.ParseUint(os.Getenv("CHANNEL_PORT"), 10, 32)
	if err != nil {
		l.Fatal("[ERROR] Unable to read port from environment.")
		return
	}

	createEventConsumers(l, byte(wid), byte(cid))

	lss := services.NewMapleSessionService(l, byte(wid), byte(cid))
	ss, err := socket.NewServer(l, lss, socket.IpAddress("0.0.0.0"), socket.Port(int(port)))
	if err != nil {
		l.Print(err.Error())
		return
	}

	registerHandlers(ss, l)
	go ss.Run()

	rs := rest.NewServer(l)
	go rs.Run()

	producers.NewChannelServer(l, context.Background()).EmitStart(byte(wid), byte(cid), ha, uint32(port))

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Println("[INFO] Shutting down via signal:", sig)
	producers.NewChannelServer(l, context.Background()).EmitShutdown(byte(wid), byte(cid), ha, uint32(port))

	sessions := registries.GetSessionRegistry().GetAll()
	for _, s := range sessions {
		s.Disconnect()
	}
}

func createEventConsumers(l *log.Logger, wid byte, cid byte) {
	createEventConsumer(l, wid, cid, "TOPIC_ENABLE_ACTIONS", consumers.EnableActionsEventCreator(), consumers.HandleEnableActionsEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHANGE_MAP_EVENT", consumers.ChangeMapEventCreator(), consumers.HandleChangeMapEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MAP_CHARACTER_EVENT", consumers.MapCharacterEventCreator(), consumers.HandleMapCharacterEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CONTROL_MONSTER_EVENT", consumers.MonsterControlEventCreator(), consumers.HandleMonsterControlEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MONSTER_EVENT", consumers.MonsterEventCreator(), consumers.HandleMonsterEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MONSTER_MOVEMENT", consumers.MonsterMovementEventCreator(), consumers.HandleMonsterMovementEvent())
}

func createEventConsumer(l *log.Logger, wid byte, cid byte, topicToken string, emptyEventCreator consumers.EmptyEventCreator, processor consumers.ChannelEventProcessor) {
	groupId := fmt.Sprintf("World Channel Coordinator %d %d", wid, cid)

	h := func(logger *log.Logger, event interface{}) {
		processor(logger, wid, cid, event)
	}

	c := consumers.NewConsumer(l, context.Background(), h,
		consumers.SetGroupId(groupId),
		consumers.SetTopicToken(topicToken),
		consumers.SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}

func registerHandlers(ss *socket.Server, l *log.Logger) {
	hr := handlerRegister(ss, l)
	hr(handler.OpCodePong, &handler.PongHandler{})
	hr(handler.OpCharacterLoggedIn, &handler.CharacterLoggedInHandler{})
	hr(handler.OpChangeMapSpecial, &handler.ChangeMapSpecialHandler{})
	hr(handler.OpMoveCharacter, &handler.MoveCharacterHandler{})
	hr(handler.OpChangeMap, &handler.ChangeMapHandler{})
	hr(handler.OpMoveLife, &handler.MoveLifeHandler{})
}

func handlerRegister(ss *socket.Server, l *log.Logger) func(uint16, request.MapleHandler) {
	return func(op uint16, handler request.MapleHandler) {
		ss.RegisterHandler(op, request.AdaptHandler(l, handler))
	}
}
