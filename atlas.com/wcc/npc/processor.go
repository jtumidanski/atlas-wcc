package npc

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

type Operator func(Model)

type SliceOperator func([]Model)

func ExecuteForEach(f Operator) SliceOperator {
	return func(npcs []Model) {
		for _, npc := range npcs {
			f(npc)
		}
	}
}

func ForEachInMap(l logrus.FieldLogger) func(mapId uint32, f Operator) {
	return func(mapId uint32, f Operator) {
		ForNPCsInMap(l)(mapId, ExecuteForEach(f))
	}
}

func ForNPCsInMap(l logrus.FieldLogger) func(mapId uint32, f SliceOperator) {
	return func(mapId uint32, f SliceOperator) {
		npcs, err := GetInMap(l)(mapId)
		if err != nil {
			return
		}
		f(npcs)
	}
}

func GetInMap(l logrus.FieldLogger) func(mapId uint32) ([]Model, error) {
	return func(mapId uint32) ([]Model, error) {
		resp, err := requestNPCsInMap(l)(mapId)
		if err != nil {
			return nil, err
		}

		ns := make([]Model, 0)
		for _, d := range resp.DataList() {
			id, err := strconv.ParseUint(d.Id, 10, 32)
			if err != nil {
				break
			}
			n := makeNPC(uint32(id), d.Attributes)
			ns = append(ns, n)
		}
		return ns, nil
	}
}

func GetInMapByObjectId(l logrus.FieldLogger) func(mapId uint32, objectId uint32) ([]Model, error) {
	return func(mapId uint32, objectId uint32) ([]Model, error) {
		resp, err := requestNPCsInMapByObjectId(l)(mapId, objectId)
		if err != nil {
			return nil, err
		}

		ns := make([]Model, 0)
		for _, d := range resp.DataList() {
			id, err := strconv.ParseUint(d.Id, 10, 32)
			if err != nil {
				break
			}
			n := makeNPC(uint32(id), d.Attributes)
			ns = append(ns, n)
		}
		return ns, nil
	}
}

func makeNPC(id uint32, att attributes) Model {
	return NewNPC(id, att.Id, att.X, att.CY, att.F, att.FH, att.RX0, att.RX1)
}
