package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/processors"
	"atlas-wcc/rest/requests"
	"atlas-wcc/socket/response/writer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type MonsterEvent struct {
	l         *log.Logger
	ctx       context.Context
	worldId   byte
	channelId byte
}

func NewMonsterEvent(l *log.Logger, ctx context.Context, worldId byte, channelId byte) *MonsterEvent {
	return &MonsterEvent{l, ctx, worldId, channelId}
}

func (mc *MonsterEvent) Init() {
	t := requests.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_MONSTER_EVENT")
	if err != nil {
		mc.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: fmt.Sprintf("World Channel Coordinator %d %d", mc.worldId, mc.channelId),
		MaxWait: 50 * time.Millisecond,
	})
	for {
		msg, err := r.ReadMessage(mc.ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		var event MonsterEventEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			mc.processEvent(event)
		}
	}
}

func (mc *MonsterEvent) processEvent(event MonsterEventEvent) {
	if mc.worldId != event.WorldId || mc.channelId != event.ChannelId {
		return
	}

	if event.Type == "CREATED" {
		mc.created(event)
	} else if event.Type == "DESTROYED" {
		mc.destroyed(event)
	}
}

func (mc *MonsterEvent) created(event MonsterEventEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		return
	}
	mc.create(*m, event)
}

func (mc *MonsterEvent) create(m domain.Monster, event MonsterEventEvent) {
	sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
	if err != nil {
		return
	}
	for _, s := range sl {
		mc.l.Printf("[INFO] spawning monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
		s.Announce(writer.WriteSpawnMonster(m, false))
	}
}

func (mc *MonsterEvent) destroyed(event MonsterEventEvent) {
	sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
	if err != nil {
		return
	}
	for _, s := range sl {
		s.Announce(writer.WriteKillMonster(event.UniqueId, false))
	}
	for _, s := range sl {
		s.Announce(writer.WriteKillMonster(event.UniqueId, true))
	}
}

type MonsterEventEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
	MonsterId uint32 `json:"monsterId"`
	ActorId   uint32 `json:"actorId"`
	Type      string `json:"type"`
}
