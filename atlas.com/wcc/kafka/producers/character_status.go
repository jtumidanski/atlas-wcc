package producers

import (
	"atlas-wcc/rest/requests"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type CharacterStatus struct {
	l   *log.Logger
	ctx context.Context
}

func NewCharacterStatus(l *log.Logger, ctx context.Context) *CharacterStatus {
	return &CharacterStatus{l, ctx}
}

func (m *CharacterStatus) EmitLogin(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	m.emit(worldId, channelId, accountId, characterId, "LOGIN")
}

func (m *CharacterStatus) EmitLogout(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	m.emit(worldId, channelId, accountId, characterId, "LOGOUT")
}

func (m *CharacterStatus) emit(worldId byte, channelId byte, accountId uint32, characterId uint32, theType string) {
	t := requests.NewTopic(m.l)
	td, err := t.GetTopic("TOPIC_CHARACTER_STATUS")
	if err != nil {
		m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        td.Attributes.Name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	e := &CharacterStatusEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		AccountId:   accountId,
		CharacterId: characterId,
		Type:        theType,
	}

	m.l.Printf("[INFO] producing CharacterStatusEvent for %d of type %s", characterId, theType)

	r, err := json.Marshal(e)
	if err != nil {
		m.l.Fatal("[ERROR] Unable to marshall event.")
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   createKey(int(characterId)),
		Value: r,
	})
	if err != nil {
		m.l.Fatal("[ERROR] Unable to produce event.")
	}
}

type CharacterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}
