package npc

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelOperator func(*Model)

type ModelListOperator func([]*Model)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r Request) ModelListProvider {
	return func(r Request) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(&v)
				if err != nil {
					return nil, err
				}
				ms = append(ms, m)
			}
			return ms, nil
		}
	}
}

func ExecuteForEach(f ModelOperator) ModelListOperator {
	return func(npcs []*Model) {
		for _, npc := range npcs {
			f(npc)
		}
	}
}

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, f ModelOperator) {
	return func(mapId uint32, f ModelOperator) {
		ForNPCsInMap(l, span)(mapId, ExecuteForEach(f))
	}
}

func ForNPCsInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, f ModelListOperator) {
	return func(mapId uint32, f ModelListOperator) {
		npcs, err := GetInMap(l, span)(mapId)
		if err != nil {
			return
		}
		f(npcs)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ModelListProvider {
	return func(mapId uint32) ModelListProvider {
		return requestModelListProvider(l, span)(requestNPCsInMap(mapId))
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]*Model, error) {
	return func(mapId uint32) ([]*Model, error) {
		return InMapModelProvider(l, span)(mapId)()
	}
}

func InMapByObjectIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, objectId uint32) ModelListProvider {
	return func(mapId uint32, objectId uint32) ModelListProvider {
		return requestModelListProvider(l, span)(requestNPCsInMapByObjectId(mapId, objectId))
	}
}

func GetInMapByObjectId(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, objectId uint32) ([]*Model, error) {
	return func(mapId uint32, objectId uint32) ([]*Model, error) {
		return InMapByObjectIdModelProvider(l, span)(mapId, objectId)()
	}
}

func makeModel(body *dataBody) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := NewNPC(uint32(id), att.Id, att.X, att.CY, att.F, att.FH, att.RX0, att.RX1)
	return &m, nil
}
