package consumers

import (
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

type EnableActions struct {
	l   *log.Logger
	ctx context.Context
}

func NewEnableActions(l *log.Logger, ctx context.Context) *EnableActions {
	return &EnableActions{l, ctx}
}

func (mc *EnableActions) Init(worldId byte, channelId byte) {
	t := requests.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_ENABLE_ACTIONS")
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

		var event EnableActionsEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			mc.processEvent(event)
		}
	}
}

func (mc *EnableActions) processEvent(event EnableActionsEvent) {
	as := getSessionForCharacterId(event.CharacterId)
	if as == nil {
		return
	}

	(*as).Announce(writer.WriteEnableActions())
}

type EnableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}
