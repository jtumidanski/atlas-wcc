package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type mpEaterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func EmptyMPEaterEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &mpEaterEvent{}
	}
}

func HandleMPEaterEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, c interface{}) {
		if event, ok := c.(*mpEaterEvent); ok {
			session.ForSessionByCharacterId(event.CharacterId, showMPEaterEffect(l, event))
			session.ForEachOtherSessionInMap(event.WorldId, event.ChannelId, event.CharacterId, showForeignMPEaterEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showMPEaterEffect(l logrus.FieldLogger, event *mpEaterEvent) session.SessionOperator {
	return func(s session.Model) {
		err := s.Announce(writer.WriteShowOwnBuff(1, event.SkillId))
		if err != nil {
			l.WithError(err).Errorf("Unable to show MP Eater application for character %d.", event.CharacterId)
		}
	}
}

func showForeignMPEaterEffect(_ logrus.FieldLogger, event *mpEaterEvent) session.SessionOperator {
	return func(s session.Model) {
		err := s.Announce(writer.WriteShowBuffEffect(event.CharacterId, 1, event.SkillId, 3))
		if err != nil {
			logrus.WithError(err).Errorf("Unable to show MP Eater effect to character %d for character %d.", s.CharacterId(), event.CharacterId)
		}
	}
}
