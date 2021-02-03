package processors

import (
	"atlas-wcc/domain"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/rest/requests"
	"strconv"
)

type MonsterOperator func(domain.Monster)

type MonstersOperator func([]domain.Monster)

func ExecuteForEachMonster(f MonsterOperator) MonstersOperator {
	return func(monsters []domain.Monster) {
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

