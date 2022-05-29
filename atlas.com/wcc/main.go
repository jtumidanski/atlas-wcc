package main

import (
	"atlas-wcc/channel"
	"atlas-wcc/character"
	"atlas-wcc/character/attack"
	"atlas-wcc/character/buff"
	"atlas-wcc/character/instruction"
	"atlas-wcc/character/inventory"
	"atlas-wcc/character/properties"
	"atlas-wcc/character/skill"
	"atlas-wcc/command"
	"atlas-wcc/configuration"
	"atlas-wcc/drop"
	"atlas-wcc/kafka"
	"atlas-wcc/logger"
	_map "atlas-wcc/map"
	"atlas-wcc/monster"
	"atlas-wcc/npc"
	"atlas-wcc/party"
	"atlas-wcc/rest"
	"atlas-wcc/server"
	"atlas-wcc/session"
	"atlas-wcc/socket"
	"atlas-wcc/tracing"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const serviceName = "atlas-wcc"
const consumerGroupFormat = "World Channel Coordinator %d %d"

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

	worldId, err := strconv.ParseUint(os.Getenv("WORLD_ID"), 10, 8)
	if err != nil {
		l.WithError(err).Fatalf("Unable to read world identifier from environment.")
		return
	}
	wid := byte(worldId)
	channelId, err := strconv.ParseUint(os.Getenv("CHANNEL_ID"), 10, 8)
	if err != nil {
		l.WithError(err).Fatalf("Unable to read channel identifier from environment.")
		return
	}
	cid := byte(channelId)
	ha := os.Getenv("HOST_ADDRESS")
	port, err := strconv.ParseUint(os.Getenv("CHANNEL_PORT"), 10, 32)
	if err != nil {
		l.WithError(err).Fatalf("Unable to read port from environment.")
		return
	}

	consumerGroupId := ConsumerGroupId(wid, cid)
	kafka.CreateConsumers(l, ctx, wg,
		buff.CharacterBuffConsumer(wid, cid)(consumerGroupId),
		buff.CharacterCancelBuffConsumer(wid, cid)(consumerGroupId),
		attack.CloseRangeConsumer(wid, cid)(consumerGroupId),
		attack.MagicAttackConsumer(wid, cid)(consumerGroupId),
		attack.RangeAttackConsumer(wid, cid)(consumerGroupId),
		character.CreatedConsumer(wid, cid)(consumerGroupId),
		_map.CharacterDamageConsumer(wid, cid)(consumerGroupId),
		character.EnableActionsConsumer(wid, cid)(consumerGroupId),
		_map.CharacterEquipmentChangedConsumer(wid, cid)(consumerGroupId),
		_map.CharacterExpressionChangeConsumer(wid, cid)(consumerGroupId),
		_map.CharacterMovementConsumer(wid, cid)(consumerGroupId),
		properties.ExperienceConsumer(wid, cid)(consumerGroupId),
		properties.MesoConsumer(wid, cid)(consumerGroupId),
		properties.StatUpdateConsumer(wid, cid)(consumerGroupId),
		drop.ReservationConsumer(wid, cid)(consumerGroupId),
		drop.PickupItemConsumer(wid, cid)(consumerGroupId),
		drop.PickupNXConsumer(wid, cid)(consumerGroupId),
		inventory.ModificationConsumer(wid, cid)(consumerGroupId),
		inventory.FullConsumer(wid, cid)(consumerGroupId),
		_map.CharacterLevelConsumer(wid, cid)(consumerGroupId),
		_map.MapCharacterConsumer(wid, cid)(consumerGroupId),
		_map.MapChangeConsumer(wid, cid)(consumerGroupId),
		_map.MessageConsumer(wid, cid)(consumerGroupId),
		_map.MPEaterConsumer(wid, cid)(consumerGroupId),
		_map.DropEventConsumer(wid, cid)(consumerGroupId),
		_map.ExpireDropConsumer(wid, cid)(consumerGroupId),
		_map.PickupConsumer(wid, cid)(consumerGroupId),
		_map.MonsterEventConsumer(wid, cid)(consumerGroupId),
		_map.MonsterMovementConsumer(wid, cid)(consumerGroupId),
		_map.MonsterDeathConsumer(wid, cid)(consumerGroupId),
		monster.ControlConsumer(wid, cid)(consumerGroupId),
		npc.TalkConsumer(wid, cid)(consumerGroupId),
		npc.TalkNumberConsumer(wid, cid)(consumerGroupId),
		npc.TalkStyleConsumer(wid, cid)(consumerGroupId),
		party.StatusConsumer(wid, cid)(consumerGroupId),
		party.MemberStatusConsumer(wid, cid)(consumerGroupId),
		_map.ReactorStatusConsumer(wid, cid)(consumerGroupId),
		server.NoticeConsumer(wid, cid)(consumerGroupId),
		skill.UpdateConsumer(wid, cid)(consumerGroupId))

	socket.CreateSocketService(l, ctx, wg)(wid, cid, int(port))

	rest.CreateService(l, ctx, wg, "/ms/csrv/worlds/{worldId}/channels/{channelId}", session.InitResource, instruction.InitResource)

	command.Registry().Add(character.AwardMesoCommandProducer(), _map.WarpMapCommandProducer())

	sl, span := tracing.StartSpan(l, "startup")
	channel.StartChannelServer(sl, span)(wid, cid, ha, uint32(port))
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
	channel.ShutdownChannelServer(sl, span)(wid, cid, ha, uint32(port))
	session.DestroyAll(sl, span, session.Registry())
	span.Finish()

	l.Infoln("Service shutdown.")
}

func ConsumerGroupId(wid byte, cid byte) string {
	return fmt.Sprintf(consumerGroupFormat, wid, cid)
}
