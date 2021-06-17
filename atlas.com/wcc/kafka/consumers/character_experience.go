package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterExperienceEvent struct {
	CharacterId  uint32 `json:"characterId"`
	PersonalGain uint32 `json:"personalGain"`
	PartyGain    uint32 `json:"partyGain"`
	Show         bool   `json:"show"`
	Chat         bool   `json:"chat"`
	White        bool   `json:"white"`
}

func CharacterExperienceEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterExperienceEvent{}
	}
}

func HandleCharacterExperienceEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterExperienceEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			if event.PersonalGain == 0 && event.PartyGain == 0 {
				return
			}

			if !event.Show {
				return
			}

			as := session.GetSessionByCharacterId(event.CharacterId)
			if as == nil {
				l.Errorf("Unable to locate session for character %d.", event.CharacterId)
				return
			}
			gain := event.PersonalGain
			party := event.PartyGain
			white := event.White
			if gain == 0 {
				gain = party
				party = 0
				white = false
			}
			err := as.Announce(writer.WriteShowExperienceGain(gain, 0, party, event.Chat, white))
			if err != nil {
				l.WithError(err).Errorf("Unable to show experience gain to character %d", as.CharacterId())
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
