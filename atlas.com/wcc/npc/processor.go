package npc

import (
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

func ForEachInMap(mapId uint32, f Operator) {
	ForNPCsInMap(mapId, ExecuteForEach(f))
}

func ForNPCsInMap(mapId uint32, f SliceOperator) {
	npcs, err := GetInMap(mapId)
	if err != nil {
		return
	}
	f(npcs)
}

func GetInMap(mapId uint32) ([]Model, error) {
	resp, err := requestNPCsInMap(mapId)
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

func GetInMapByObjectId(mapId uint32, objectId uint32) ([]Model, error) {
	resp, err := requestNPCsInMapByObjectId(mapId, objectId)
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

func makeNPC(id uint32, att attributes) Model {
	return NewNPC(id, att.Id, att.X, att.CY, att.F, att.FH, att.RX0, att.RX1)
}
