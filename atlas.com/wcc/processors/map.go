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

type MonsterOperator func(domain.Monster)

type MonstersOperator func([]domain.Monster)

func ExecuteForEachMonster(f MonsterOperator) MonstersOperator {
	return func(monsters []domain.Monster) {
		for _, monster := range monsters {
			f(monster)
		}
	}
}

type DropOperator func(domain.Drop)

type DropsOperator func([]domain.Drop)

func ExecuteForEachDrop(f DropOperator) DropsOperator {
	return func(drops []domain.Drop) {
		for _, drop := range drops {
			f(drop)
		}
	}
}

func GetCharacterIdsInMap(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	resp, err := requests.MapRegistry().GetCharactersInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	cIds := make([]uint32, 0)
	for _, d := range resp.DataList() {
		cId, err := strconv.ParseUint(d.Id, 10, 32)
		if err != nil {
			break
		}
		cIds = append(cIds, uint32(cId))
	}
	return cIds, nil
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

func ForEachMonsterInMap(worldId byte, channelId byte, mapId uint32, f MonsterOperator) {
	ForMonstersInMap(worldId, channelId, mapId, ExecuteForEachMonster(f))
}

func ForMonstersInMap(worldId byte, channelId byte, mapId uint32, f MonstersOperator) {
	monsters, err := GetMonstersInMap(worldId, channelId, mapId)
	if err != nil {
		return
	}
	f(monsters)
}

func GetMonstersInMap(worldId byte, channelId byte, mapId uint32) ([]domain.Monster, error) {
	resp, err := requests.MonsterRegistry().GetInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	ns := make([]domain.Monster, 0)
	for _, d := range resp.DataList() {
		id, err := strconv.ParseUint(d.Id, 10, 32)
		if err != nil {
			break
		}
		n := makeMonster(uint32(id), d.Attributes)
		ns = append(ns, n)
	}
	return ns, nil
}

func GetMonster(id uint32) (*domain.Monster, error) {
	resp, err := requests.MonsterRegistry().GetById(id)
	if err != nil {
		return nil, err
	}

	d := resp.Data()
	n := makeMonster(id, d.Attributes)
	return &n, nil
}

func makeMonster(id uint32, att attributes.MonsterAttributes) domain.Monster {
	return domain.NewMonster(id, att.ControlCharacterId, att.MonsterId, att.X, att.Y, att.Stance, att.FH, att.Team)
}

func ForEachDropInMap(worldId byte, channelId byte, mapId uint32, f DropOperator) {
	ForDropsInMap(worldId, channelId, mapId, ExecuteForEachDrop(f))
}

func ForDropsInMap(worldId byte, channelId byte, mapId uint32, f DropsOperator) {
	drops, err := GetDropsInMap(worldId, channelId, mapId)
	if err != nil {
		return
	}
	f(drops)
}

func GetDropsInMap(worldId byte, channelId byte, mapId uint32) ([]domain.Drop, error) {
	resp, err := requests.DropRegistry().GetDropsInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	ns := make([]domain.Drop, 0)
	for _, d := range resp.DataList() {
		id, err := strconv.ParseUint(d.Id, 10, 32)
		if err != nil {
			break
		}
		n := makeDrop(uint32(id), d.Attributes)
		ns = append(ns, n)
	}
	return ns, nil
}

func makeDrop(id uint32, att attributes.DropAttributes) domain.Drop {
	return domain.NewDropBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetMapId(att.MapId).
		SetUniqueId(id).
		SetItemId(att.ItemId).
		SetMeso(att.Meso).
		SetDropType(att.DropType).
		SetDropX(att.DropX).
		SetDropY(att.DropY).
		SetOwnerId(att.OwnerId).
		SetOwnerPartyId(att.OwnerPartyId).
		SetDropTime(att.DropTime).
		SetDropperUniqueId(att.DropperUniqueId).
		SetDropperX(att.DropperX).
		SetDropperY(att.DropperY).
		SetCharacterDrop(att.CharacterDrop).
		SetMod(att.Mod).
		Build()
}
