package processors

import (
	"atlas-wcc/domain"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/rest/requests"
	"strconv"
)
type NPCOperator func(domain.NPC)

type NPCsOperator func([]domain.NPC)

func ExecuteForEachNPC(f NPCOperator) NPCsOperator {
	return func(npcs []domain.NPC) {
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

func GetNPCsInMap(mapId uint32) ([]domain.NPC, error) {
	resp, err := requests.MapInformation().GetNPCsInMap(mapId)
	if err != nil {
		return nil, err
	}

	ns := make([]domain.NPC, 0)
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

func GetNPCsInMapByObjectId(mapId uint32, objectId uint32) ([]domain.NPC, error) {
	resp, err := requests.MapInformation().GetNPCsInMapByObjectId(mapId, objectId)
	if err != nil {
		return nil, err
	}

	ns := make([]domain.NPC, 0)
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

func makeNPC(id uint32, att attributes.NpcAttributes) domain.NPC {
	return domain.NewNPC(id, att.Id, att.X, att.CY, att.F, att.FH, att.RX0, att.RX1)
}