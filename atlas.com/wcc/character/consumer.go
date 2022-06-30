package character

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/kafka"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameCharacterCreated = "character_created_event"
	consumerNameEnableActions    = "enable_actions_command"
	topicTokenCreated            = "TOPIC_CHARACTER_CREATED_EVENT"
	topicTokenEnableActions      = "TOPIC_ENABLE_ACTIONS"
)

func CreatedConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[createdEvent](consumerNameCharacterCreated, topicTokenCreated, groupId, handleCreated(wid, cid))
	}
}

const characterCreatedFormat = "Character %s has been created."

type createdEvent struct {
	WorldId     byte   `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Name        string `json:"name"`
}

func handleCreated(wid byte, _ byte) kafka.HandlerFunc[createdEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event createdEvent) {
		if wid != event.WorldId {
			return
		}

		session.ForEachGM(session.Announce(writer.WriteYellowTip(l)(fmt.Sprintf(characterCreatedFormat, event.Name))))
	}
}

func EnableActionsConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[enableActionsEvent](consumerNameEnableActions, topicTokenEnableActions, groupId, handleEnableActions(wid, cid))
	}
}

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func handleEnableActions(_ byte, _ byte) kafka.HandlerFunc[enableActionsEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event enableActionsEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, session.Announce(properties.WriteEnableActions(l)))
	}
}
