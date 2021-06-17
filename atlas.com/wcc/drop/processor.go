package drop

import (
	"strconv"
)

type Operator func(Model)

type SliceOperator func([]Model)

func ExecuteForEach(f Operator) SliceOperator {
	return func(drops []Model) {
		for _, drop := range drops {
			f(drop)
		}
	}
}

func ForEachInMap(worldId byte, channelId byte, mapId uint32, f Operator) {
	ForDropsInMap(worldId, channelId, mapId, ExecuteForEach(f))
}

func ForDropsInMap(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	drops, err := GetInMap(worldId, channelId, mapId)
	if err != nil {
		return
	}
	f(drops)
}

func GetInMap(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	resp, err := requestDropsInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	ns := make([]Model, 0)
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

func makeDrop(id uint32, att attributes) Model {
	return NewDropBuilder().
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
