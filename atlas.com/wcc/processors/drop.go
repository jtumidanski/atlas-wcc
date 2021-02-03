package processors

import (
	"atlas-wcc/domain"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/rest/requests"
	"strconv"
)

type DropOperator func(domain.Drop)

type DropsOperator func([]domain.Drop)

func ExecuteForEachDrop(f DropOperator) DropsOperator {
	return func(drops []domain.Drop) {
		for _, drop := range drops {
			f(drop)
		}
	}
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