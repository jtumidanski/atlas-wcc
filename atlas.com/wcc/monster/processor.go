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
		model.ForEach(InMapModelProvider(l, span)(worldId, channelId, mapId), f)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[Model] {
	return func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(worldId, channelId, mapId), makeModel)
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

func SpawnSessionOperator(l logrus.FieldLogger) func(monster Model) model.Operator[session.Model] {
	return func(m Model) model.Operator[session.Model] {
		return session.AnnounceOperator(WriteSpawnMonster(l)(m, false))
	}
}

func SpawnForSession(l logrus.FieldLogger) func(s session.Model) model.Operator[Model] {
	return func(s session.Model) model.Operator[Model] {
		return func(m Model) error {
			return SpawnSessionOperator(l)(m)(s)
		}
	}
}
