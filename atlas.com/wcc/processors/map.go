package processors

import (
	"atlas-wcc/domain"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/rest/requests"
	"strconv"
)

func GetCharacterIdsInMap(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	resp, err := requests.GetCharactersInMap(worldId, channelId, mapId)
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

func GetNPCsInMap(mapId uint32) ([]domain.NPC, error) {
	resp, err := requests.GetNPCsInMap(mapId)
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

func GetMonstersInMap(worldId byte, channelId byte, mapId uint32) ([]domain.Monster, error) {
	resp, err := requests.GetMonstersInMap(worldId, channelId, mapId)
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
	resp, err := requests.GetMonster(id)
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

func GetDropsInMap(worldId byte, channelId byte, mapId uint32) ([]domain.Drop, error) {
	resp, err := requests.GetDropsInMap(worldId, channelId, mapId)
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