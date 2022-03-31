package npc

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, f model.Operator[Model]) {
	return func(mapId uint32, f model.Operator[Model]) {
		ForNPCsInMap(l, span)(mapId, model.ExecuteForEach(f))
	}
}

func ForNPCsInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, f model.SliceOperator[Model]) {
	return func(mapId uint32, f model.SliceOperator[Model]) {
		npcs, err := GetInMap(l, span)(mapId)
		if err != nil {
			return
		}
		f(npcs)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) model.SliceProvider[Model] {
	return func(mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestNPCsInMap(mapId), makeModel)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]Model, error) {
	return func(mapId uint32) ([]Model, error) {
		return InMapModelProvider(l, span)(mapId)()
	}
}

func InMapByObjectIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, objectId uint32) model.SliceProvider[Model] {
	return func(mapId uint32, objectId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestNPCsInMapByObjectId(mapId, objectId), makeModel)
	}
}

func GetInMapByObjectId(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, objectId uint32) ([]Model, error) {
	return func(mapId uint32, objectId uint32) ([]Model, error) {
		return InMapByObjectIdModelProvider(l, span)(mapId, objectId)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := body.Attributes
	m := NewNPC(uint32(id), att.Id, att.X, att.CY, att.F, att.FH, att.RX0, att.RX1)
	return m, nil
}

func SpawnNPCForSession(l logrus.FieldLogger) func(s session.Model) model.Operator[Model] {
	return func(s session.Model) model.Operator[Model] {
		return func(n Model) {
			err := session.Announce(WriteSpawnNPC(l)(n))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn npc %d for character %d", n.Id(), s.CharacterId())
			}
			err = session.Announce(WriteSpawnNPCController(l)(n, true))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn npc controller %d for character %d", n.Id(), s.CharacterId())
			}
		}
	}
}
