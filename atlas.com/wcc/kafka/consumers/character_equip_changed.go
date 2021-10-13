package consumers

import (
	"atlas-wcc/character"
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
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
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterEquipChangedEvent); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachOtherInMap(l, span)(wid, cid, event.CharacterId, updateCharacterLook(l, span)(event.CharacterId))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func updateCharacterLook(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) session.Operator {
	return func(characterId uint32) session.Operator {
		return func(s *session.Model) {
			r, err := character.GetCharacterById(l, span)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve character %d details.", s.CharacterId())
				return
			}
			c, err := character.GetCharacterById(l, span)(characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve character %d details.", characterId)
				return
			}
			err = s.Announce(writer.WriteCharacterLookUpdated(l)(*r, *c))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to %d that character %d has changed their look.", s.CharacterId(), characterId)
			}
		}
	}
}
