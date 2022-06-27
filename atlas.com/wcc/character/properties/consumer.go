package properties

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameCharacterExperience = "character_experience_event"
	consumerNameMesoGained          = "meso_gained_event"
	consumerNameStatistic           = "character_statistic_event"
	topicTokenExperienceGained      = "TOPIC_CHARACTER_EXPERIENCE_EVENT"
	topicTokenMeso                  = "TOPIC_MESO_GAINED"
	topicTokenStatistic             = "TOPIC_CHARACTER_STAT_EVENT"
)

func ExperienceConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[experienceEvent](consumerNameCharacterExperience, topicTokenExperienceGained, groupId, handleExperienceGain(wid, cid))
	}
}

type experienceEvent struct {
	CharacterId  uint32 `json:"characterId"`
	PersonalGain uint32 `json:"personalGain"`
	PartyGain    uint32 `json:"partyGain"`
	Show         bool   `json:"show"`
	Chat         bool   `json:"chat"`
	White        bool   `json:"white"`
}

func handleExperienceGain(_ byte, _ byte) kafka.HandlerFunc[experienceEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event experienceEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		if event.PersonalGain == 0 && event.PartyGain == 0 {
			return
		}

		if !event.Show {
			return
		}

		as, err := session.GetByCharacterId(event.CharacterId)
		if err != nil {
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
		err = session.Announce(WriteShowExperienceGain(l)(gain, 0, party, event.Chat, white))(as)
		if err != nil {
			l.WithError(err).Errorf("Unable to show experience gain to character %d", as.CharacterId())
		}
	}
}

func MesoConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[mesoEvent](consumerNameMesoGained, topicTokenMeso, groupId, handleMeso(wid, cid))
	}
}

type mesoEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        int32  `json:"gain"`
}

func handleMeso(_ byte, _ byte) kafka.HandlerFunc[mesoEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event mesoEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, showChange(l, event))
	}
}

func showChange(l logrus.FieldLogger, event mesoEvent) model.Operator[session.Model] {
	mg := WriteShowMesoGain(l)(event.Gain, false)
	ea := WriteEnableActions(l)
	return func(s session.Model) error {
		err := session.Announce(mg)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			return err
		}
		err = session.Announce(ea)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		return err
	}
}

func StatUpdateConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[statisticEvent](consumerNameStatistic, topicTokenStatistic, groupId, handleStatisticChange(wid, cid))
	}
}

type statisticEvent struct {
	CharacterId uint32   `json:"characterId"`
	Updates     []string `json:"updates"`
}

func handleStatisticChange(_ byte, _ byte) kafka.HandlerFunc[statisticEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statisticEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, updateStats(l, span, event))
	}
}

func updateStats(l logrus.FieldLogger, span opentracing.Span, event statisticEvent) model.Operator[session.Model] {
	return func(s session.Model) error {
		ca, err := GetById(l, span)(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrive character %d properties", event.CharacterId)
			return err
		}

		var statUpdates []StatUpdate
		for _, t := range event.Updates {
			statUpdates = append(statUpdates, getStatUpdate(ca, t))
		}
		err = session.Announce(WriteCharacterStatUpdate(l)(statUpdates, true))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write character stat update for %d", event.CharacterId)
		}
		return err
	}
}

func getStatUpdate(ca Model, stat string) StatUpdate {
	switch stat {
	case "EXPERIENCE":
		return NewStatUpdate(StatUpdateExperience, ca.Experience())
	case "SKIN":
		return NewStatUpdate(StatUpdateSkin, uint32(ca.SkinColor()))
	case "FACE":
		return NewStatUpdate(StatUpdateFace, ca.Face())
	case "HAIR":
		return NewStatUpdate(StatUpdateHair, ca.Hair())
	case "LEVEL":
		return NewStatUpdate(StatUpdateLevel, uint32(ca.Level()))
	case "AVAILABLE_AP":
		return NewStatUpdate(StatUpdateAvailableAP, uint32(ca.Ap()))
	case "AVAILABLE_SP":
		return NewStatUpdate(StatUpdateAvailableSP, uint32(ca.Sp()[0]))
	case "HP":
		return NewStatUpdate(StatUpdateHP, uint32(ca.Hp()))
	case "MP":
		return NewStatUpdate(StatUpdateMP, uint32(ca.Mp()))
	case "MAX_HP":
		return NewStatUpdate(StatUpdateMaxHP, uint32(ca.MaxHp()))
	case "MAX_MP":
		return NewStatUpdate(StatUpdateMaxMP, uint32(ca.MaxMp()))
	case "JOB":
		return NewStatUpdate(StatUpdateJob, uint32(ca.JobId()))
	case "STRENGTH":
		return NewStatUpdate(StatUpdateStrength, uint32(ca.Strength()))
	case "DEXTERITY":
		return NewStatUpdate(StatUpdateDexterity, uint32(ca.Dexterity()))
	case "LUCK":
		return NewStatUpdate(StatUpdateLuck, uint32(ca.Luck()))
	case "INTELLIGENCE":
		return NewStatUpdate(StatUpdateIntelligence, uint32(ca.Intelligence()))
	case "FAME":
		return NewStatUpdate(StatUpdateFame, uint32(ca.Fame()))
	case "MESO":
		return NewStatUpdate(StatUpdateMeso, ca.Meso())
	case "PET":
		return NewStatUpdate(StatUpdatePet, 0)
	case "GACHAPON_EXPERIENCE":
		return NewStatUpdate(StatUpdateGachaponExperience, ca.GachaponExperience())
	}
	panic("unknown stat update type")
}
