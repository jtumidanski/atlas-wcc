package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type CharacterStatUpdateEvent struct {
	CharacterId uint32   `json:"characterId"`
	Updates     []string `json:"updates"`
}

func CharacterStatUpdateEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &CharacterStatUpdateEvent{}
	}
}

func HandleCharacterStatUpdateEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*CharacterStatUpdateEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, updateStats(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterStatUpdateEvent]")
		}
	}
}

func updateStats(_ *log.Logger, event *CharacterStatUpdateEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		ca, err := processors.GetCharacterAttributesById(event.CharacterId)
		if err != nil {
			return
		}

		var statUpdates []writer.StatUpdate
		for _, t := range event.Updates {
			statUpdates = append(statUpdates, getStatUpdate(ca, t))
		}
		session.Announce(writer.WriteCharacterStatUpdate(statUpdates, true))
	}
}

func getStatUpdate(ca *domain.CharacterAttributes, stat string) writer.StatUpdate {
	switch stat {
	case "EXPERIENCE":
		return writer.NewStatUpdate(writer.StatUpdateExperience, ca.Experience())
	case "SKIN":
		return writer.NewStatUpdate(writer.StatUpdateSkin, uint32(ca.SkinColor()))
	case "FACE":
		return writer.NewStatUpdate(writer.StatUpdateFace, ca.Face())
	case "HAIR":
		return writer.NewStatUpdate(writer.StatUpdateHair, ca.Hair())
	case "LEVEL":
		return writer.NewStatUpdate(writer.StatUpdateLevel, uint32(ca.Level()))
	case "AVAILABLE_AP":
		return writer.NewStatUpdate(writer.StatUpdateAvailableAP, uint32(ca.Ap()))
	case "AVAILABLE_SP":
		return writer.NewStatUpdate(writer.StatUpdateAvailableSP, 0)
	case "HP":
		return writer.NewStatUpdate(writer.StatUpdateHP, uint32(ca.Hp()))
	case "MP":
		return writer.NewStatUpdate(writer.StatUpdateMP, uint32(ca.Mp()))
	case "MAX_HP":
		return writer.NewStatUpdate(writer.StatUpdateMaxHP, uint32(ca.MaxHp()))
	case "MAX_MP":
		return writer.NewStatUpdate(writer.StatUpdateMaxMP, uint32(ca.MaxMp()))
	case "JOB":
		return writer.NewStatUpdate(writer.StatUpdateJob, uint32(ca.JobId()))
	case "STRENGTH":
		return writer.NewStatUpdate(writer.StatUpdateStrength, uint32(ca.Strength()))
	case "DEXTERITY":
		return writer.NewStatUpdate(writer.StatUpdateDexterity, uint32(ca.Dexterity()))
	case "LUCK":
		return writer.NewStatUpdate(writer.StatUpdateLuck, uint32(ca.Luck()))
	case "INTELLIGENCE":
		return writer.NewStatUpdate(writer.StatUpdateIntelligence, uint32(ca.Intelligence()))
	case "FAME":
		return writer.NewStatUpdate(writer.StatUpdateFame, uint32(ca.Fame()))
	case "MESO":
		return writer.NewStatUpdate(writer.StatUpdateMeso, ca.Meso())
	case "PET":
		return writer.NewStatUpdate(writer.StatUpdatePet, 0)
	case "GACHAPON_EXPERIENCE":
		return writer.NewStatUpdate(writer.StatUpdateGachaponExperience, ca.GachaponExperience())
	}
	panic("unknown stat update type")
}
