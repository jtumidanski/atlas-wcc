package consumers

import (
	"atlas-wcc/character"
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterEquipChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Change      string `json:"change"`
}

func CharacterEquipChangedEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterEquipChangedEvent{}
	}
}

func HandleCharacterEquipChangedEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterEquipChangedEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachOtherSessionInMap(wid, cid, event.CharacterId, updateCharacterLook(l, event.CharacterId))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func updateCharacterLook(l logrus.FieldLogger, characterId uint32) session.SessionOperator {
	return func(s session.Model) {
		r, err := character.GetCharacterById(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d details.", s.CharacterId())
			return
		}
		c, err := character.GetCharacterById(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d details.", characterId)
			return
		}
		err = s.Announce(writer.WriteCharacterLookUpdated(*r, *c))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to %d that character %d has changed their look.", s.CharacterId(), characterId)
		}
	}
}
