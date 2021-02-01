package consumers

import (
	"atlas-wcc/mapleSession"
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

type MonsterControl struct {
	l   *log.Logger
	ctx context.Context
}

func NewMonsterControl(l *log.Logger, ctx context.Context) *MonsterControl {
	return &MonsterControl{l, ctx}
}

func (mc *MonsterControl) Init(worldId byte, channelId byte) {
	t := requests.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_CONTROL_MONSTER_EVENT")
	if err != nil {
		mc.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: fmt.Sprintf("World Channel Coordinator %d %d", worldId, channelId),
		MaxWait: 50 * time.Millisecond,
	})
	for {
		msg, err := r.ReadMessage(mc.ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		var event MonsterControlEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			mc.processEvent(event)
		}
	}
}

func (mc *MonsterControl) processEvent(event MonsterControlEvent) {
	as := getSessionForCharacterId(event.CharacterId)
	if as == nil {
		mc.l.Printf("[ERROR] cannot location session for character %d for monster control event processing.", event.CharacterId)
		return
	}

	if event.Type == "START" {
		mc.start(*as, event)
	} else if event.Type == "STOP" {
		mc.stop(*as, event)
	}
}

func (mc *MonsterControl) start(s mapleSession.MapleSession, event MonsterControlEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		mc.l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
		return
	}
	mc.l.Printf("[INFO] controlling monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
	s.Announce(writer.WriteControlMonster(m, false, false))
}

func (mc *MonsterControl) stop(s mapleSession.MapleSession, event MonsterControlEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		mc.l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
		return
	}
	mc.l.Printf("[INFO] removing control of monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
	s.Announce(writer.WriteStopControlMonster(m))
}

type MonsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    uint32 `json:"uniqueId"`
	Type        string `json:"type"`
}
