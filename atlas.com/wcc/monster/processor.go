package monster

import (
	"strconv"
)

type MonsterOperator func(Model)

type MonstersOperator func([]Model)

func ExecuteForEachMonster(f MonsterOperator) MonstersOperator {
	return func(monsters []Model) {
		for _, monster := range monsters {
			f(monster)
		}
	}
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

func GetMonstersInMap(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	resp, err := MonsterRegistry().GetInMap(worldId, channelId, mapId)
	if err != nil {
		return nil, err
	}

	ns := make([]Model, 0)
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

func GetMonster(id uint32) (*Model, error) {
	resp, err := MonsterRegistry().GetById(id)
	if err != nil {
		return nil, err
	}

	d := resp.Data()
	n := makeMonster(id, d.Attributes)
	return &n, nil
}

func makeMonster(id uint32, att MonsterAttributes) Model {
	return NewMonster(id, att.ControlCharacterId, att.MonsterId, att.X, att.Y, att.Stance, att.FH, att.Team)
}

