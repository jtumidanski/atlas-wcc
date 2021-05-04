package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterEquipChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Change      string `json:"change"`
}

func CharacterEquipChangedEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterEquipChangedEvent{}
	}
}

func HandleCharacterEquipChangedEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterEquipChangedEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachOtherSessionInMap(wid, cid, event.CharacterId, updateCharacterLook(l, event.CharacterId))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func updateCharacterLook(l logrus.FieldLogger, characterId uint32) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		r, err := processors.GetCharacterById(session.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d details.", session.CharacterId())
			return
		}
		c, err := processors.GetCharacterById(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d details.", characterId)
			return
		}
		err = session.Announce(writer.WriteCharacterLookUpdated(*r, *c))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to %d that character %d has changed their look.", session.CharacterId(), characterId)
		}
	}
}
