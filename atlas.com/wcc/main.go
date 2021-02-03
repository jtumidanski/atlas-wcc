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

const (
	consumerGroupFormat = "World Channel Coordinator %d %d"
)

func main() {
	l := log.New(os.Stdout, "wcc ", log.LstdFlags|log.Lmicroseconds)

	_, err := registries.GetConfiguration()
	if err != nil {
		l.Fatal("[ERROR] unable to successfully load configuration.")
	}

	wid, err := strconv.ParseUint(os.Getenv("WORLD_ID"), 10, 8)
	if err != nil {
		l.Fatal("[ERROR] unable to read world identifier from environment.")
		return
	}
	cid, err := strconv.ParseUint(os.Getenv("CHANNEL_ID"), 10, 8)
	if err != nil {
		l.Fatal("[ERROR] unable to read channel identifier from environment.")
		return
	}
	ha := os.Getenv("HOST_ADDRESS")
	port, err := strconv.ParseUint(os.Getenv("CHANNEL_PORT"), 10, 32)
	if err != nil {
		l.Fatal("[ERROR] unable to read port from environment.")
		return
	}

	createEventConsumers(l, byte(wid), byte(cid))
	createSocketService(l, wid, cid, err, port)
	createRestService(l)

	producers.ChannelServer(l, context.Background()).Start(byte(wid), byte(cid), ha, uint32(port))

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Println("[INFO] shutting down via signal:", sig)
	producers.ChannelServer(l, context.Background()).Shutdown(byte(wid), byte(cid), ha, uint32(port))

	sessions := registries.GetSessionRegistry().GetAll()
	for _, s := range sessions {
		s.Disconnect()
	}
}

func createRestService(l *log.Logger) {
	rs := rest.NewServer(l)
	go rs.Run()
}

func createSocketService(l *log.Logger, wid uint64, cid uint64, err error, port uint64) {
	lss := services.NewMapleSessionService(l, byte(wid), byte(cid))
	ss, err := socket.NewServer(l, lss, socket.IpAddress("0.0.0.0"), socket.Port(int(port)))
	if err != nil {
		l.Fatal(err.Error())
	}

	registerSocketRequestHandlers(ss, l)
	go ss.Run()
}

func createEventConsumers(l *log.Logger, wid byte, cid byte) {
	createEventConsumer(l, wid, cid, "TOPIC_ENABLE_ACTIONS", consumers.EnableActionsEventCreator(), consumers.HandleEnableActionsEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHANGE_MAP_EVENT", consumers.ChangeMapEventCreator(), consumers.HandleChangeMapEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MAP_CHARACTER_EVENT", consumers.MapCharacterEventCreator(), consumers.HandleMapCharacterEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CONTROL_MONSTER_EVENT", consumers.MonsterControlEventCreator(), consumers.HandleMonsterControlEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MONSTER_EVENT", consumers.MonsterEventCreator(), consumers.HandleMonsterEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MONSTER_MOVEMENT", consumers.MonsterMovementEventCreator(), consumers.HandleMonsterMovementEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_MOVEMENT", consumers.CharacterMovementEventCreator(), consumers.HandleCharacterMovementEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_MAP_MESSAGE_EVENT", consumers.CharacterMapMessageEventCreator(), consumers.HandleCharacterMapMessageEvent())
	createEventConsumer(l, wid, cid, "EXPRESSION_CHANGED", consumers.CharacterExpressionChangedEventCreator(), consumers.HandleCharacterExpressionChangedEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_CREATED_EVENT", consumers.CharacterCreatedEventCreator(), consumers.HandleCharacterCreatedEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_EXPERIENCE_EVENT", consumers.CharacterExperienceEventCreator(), consumers.HandleCharacterExperienceEvent())
	createEventConsumer(l, wid, cid, "TOPIC_INVENTORY_MODIFICATION", consumers.CharacterInventoryModificationEventCreator(), consumers.HandleCharacterInventoryModificationEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_LEVEL_EVENT", consumers.CharacterLevelEventCreator(), consumers.HandleCharacterLevelEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MESO_GAINED", consumers.CharacterMesoEventCreator(), consumers.HandleCharacterMesoEvent())
	createEventConsumer(l, wid, cid, "TOPIC_PICKED_UP_ITEM", consumers.ItemPickedUpEventCreator(), consumers.HandleItemPickedUpEvent())
	createEventConsumer(l, wid, cid, "TOPIC_PICKED_UP_NX", consumers.NXPickedUpEventCreator(), consumers.HandleNXPickedUpEvent())
	createEventConsumer(l, wid, cid, "TOPIC_DROP_RESERVATION_EVENT", consumers.DropReservationEventCreator(), consumers.HandleDropReservationEvent())
	createEventConsumer(l, wid, cid, "TOPIC_PICKUP_DROP_EVENT", consumers.DropPickedUpEventCreator(), consumers.HandleDropPickedUpEvent())
	createEventConsumer(l, wid, cid, "TOPIC_DROP_EVENT", consumers.DropEventCreator(), consumers.HandleDropEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_SKILL_UPDATE_EVENT", consumers.CharacterSkillUpdateEventCreator(), consumers.HandleCharacterSkillUpdateEvent())
	createEventConsumer(l, wid, cid, "TOPIC_CHARACTER_STAT_EVENT", consumers.CharacterStatUpdateEventCreator(), consumers.HandleCharacterStatUpdateEvent())
	createEventConsumer(l, wid, cid, "TOPIC_SERVER_NOTICE_COMMAND", consumers.ServerNoticeEventCreator(), consumers.HandleServerNoticeEvent())
	createEventConsumer(l, wid, cid, "TOPIC_MONSTER_KILLED_EVENT", consumers.MonsterKilledEventCreator(), consumers.HandleMonsterKilledEvent())
	createEventConsumer(l, wid, cid, "TOPIC_SHOW_DAMAGE_CHARACTER_COMMAND", consumers.CharacterDamagedEventCreator(), consumers.HandleCharacterDamagedEvent())
	createEventConsumer(l, wid, cid, "TOPIC_NPC_TALK_COMMAND", consumers.NPCTalkEventCreator(), consumers.HandleNPCTalkEvent())
}

func createEventConsumer(l *log.Logger, wid byte, cid byte, topicToken string, emptyEventCreator consumers.EmptyEventCreator, processor consumers.ChannelEventProcessor) {
	groupId := fmt.Sprintf(consumerGroupFormat, wid, cid)

	h := func(logger *log.Logger, event interface{}) {
		processor(logger, wid, cid, event)
	}

	c := consumers.NewConsumer(l, context.Background(), h,
		consumers.SetGroupId(groupId),
		consumers.SetTopicToken(topicToken),
		consumers.SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}

func registerSocketRequestHandlers(ss *socket.Server, l *log.Logger) {
	hr := socketRequestHandlerRegistration(ss, l)
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

func socketRequestHandlerRegistration(ss *socket.Server, l *log.Logger) func(uint16, request.SessionStateValidator, request.SessionRequestHandler) {
	return func(op uint16, validator request.SessionStateValidator, handler request.SessionRequestHandler) {
		ss.RegisterHandler(op, request.AdaptHandler(l, validator, handler))
	}
}
