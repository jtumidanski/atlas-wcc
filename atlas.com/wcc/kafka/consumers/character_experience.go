package consumers

import (
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type characterExperienceEvent struct {
	CharacterId  uint32 `json:"characterId"`
	PersonalGain uint32 `json:"personalGain"`
	PartyGain    uint32 `json:"partyGain"`
	Show         bool   `json:"show"`
	Chat         bool   `json:"chat"`
	White        bool   `json:"white"`
}

func CharacterExperienceEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterExperienceEvent{}
	}
}

func HandleCharacterExperienceEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterExperienceEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			if event.PersonalGain == 0 && event.PartyGain == 0 {
				return
			}

			if !event.Show {
				return
			}

			as := processors.GetSessionByCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
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
			(*as).Announce(writer.WriteShowExperienceGain(gain, 0, party, event.Chat, white))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterExperienceEvent]")
		}
	}
}
