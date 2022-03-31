package monster

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f model.Operator[Model]) {
	return func(worldId byte, channelId byte, mapId uint32, f model.Operator[Model]) {
		ForInMap(l, span)(worldId, channelId, mapId, model.ExecuteForEach(f))
	}
}

func ForInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f model.SliceOperator[Model]) {
	return func(worldId byte, channelId byte, mapId uint32, f model.SliceOperator[Model]) {
		monsters, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(monsters)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[Model] {
	return func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(worldId, channelId, mapId), makeModel)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
		return InMapModelProvider(l, span)(worldId, channelId, mapId)()
	}
}

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(id uint32) model.Provider[Model] {
	return func(id uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(id), makeModel)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(id uint32) (Model, error) {
	return func(id uint32) (Model, error) {
		return ByIdModelProvider(l, span)(id)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := body.Attributes
	m := NewMonster(uint32(id), att.ControlCharacterId, att.MonsterId, att.X, att.Y, att.Stance, att.FH, att.Team)
	return m, nil
}

func SpawnForSession(l logrus.FieldLogger) func(s session.Model) model.Operator[Model] {
	return func(s session.Model) model.Operator[Model] {
		return func(m Model) {
			err := session.Announce(WriteSpawnMonster(l)(m, false))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn monster %d for character %d", m.MonsterId(), s.CharacterId())
			}
		}
	}
}

func DestroyForSession(l logrus.FieldLogger, uniqueId uint32) model.Operator[session.Model] {
	k1 := WriteKillMonster(l)(uniqueId, false)
	k2 := WriteKillMonster(l)(uniqueId, true)
	return func(s session.Model) {
		err := session.Announce(k1)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = session.Announce(k2)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func CreateForSession(l logrus.FieldLogger, m Model) model.Operator[session.Model] {
	sm := WriteSpawnMonster(l)(m, false)
	return func(s session.Model) {
		err := session.Announce(sm)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func KillForSession(l logrus.FieldLogger, uniqueId uint32) model.Operator[session.Model] {
	b := WriteKillMonster(l)(uniqueId, true)
	return func(s session.Model) {
		err := session.Announce(b)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func MoveForSession(l logrus.FieldLogger, objectId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, movementList []byte) model.Operator[session.Model] {
	b := WriteMoveMonster(l)(objectId, skillPossible, skill, skillId, skillLevel, option, startX, startY, movementList)
	return func(s session.Model) {
		err := session.Announce(b)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
