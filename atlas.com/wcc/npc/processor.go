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
		model.ForEach(InMapModelProvider(l, span)(mapId), f)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) model.SliceProvider[Model] {
	return func(mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestNPCsInMap(mapId), makeModel)
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

func SpawnSessionOperator(l logrus.FieldLogger) func(npc Model) model.Operator[session.Model] {
	return func(npc Model) model.Operator[session.Model] {
		return session.AnnounceOperator(WriteSpawnNPC(l)(npc), WriteSpawnNPCController(l)(npc, true))
	}
}

func SpawnNPCForSession(l logrus.FieldLogger) func(s session.Model) model.Operator[Model] {
	return func(s session.Model) model.Operator[Model] {
		return func(n Model) error {
			return SpawnSessionOperator(l)(n)(s)
		}
	}
}
