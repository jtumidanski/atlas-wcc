package npc

import (
	"strconv"
)

type NPCOperator func(Model)

type NPCsOperator func([]Model)

func ExecuteForEachNPC(f NPCOperator) NPCsOperator {
	return func(npcs []Model) {
		for _, npc := range npcs {
			f(npc)
		}
	}
}

func ForEachNPCInMap(mapId uint32, f NPCOperator) {
	ForNPCsInMap(mapId, ExecuteForEachNPC(f))
}

func ForNPCsInMap(mapId uint32, f NPCsOperator) {
	npcs, err := GetNPCsInMap(mapId)
	if err != nil {
		return
	}
	f(npcs)
}

func GetNPCsInMap(mapId uint32) ([]Model, error) {
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

func GetNPCsInMapByObjectId(mapId uint32, objectId uint32) ([]Model, error) {
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
