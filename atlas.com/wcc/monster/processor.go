package monster

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

type Operator func(Model)

type SliceOperator func([]Model)

func ExecuteForEach(f Operator) SliceOperator {
	return func(monsters []Model) {
		for _, monster := range monsters {
			f(monster)
		}
	}
}

func ForEachInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, f Operator) {
	return func(worldId byte, channelId byte, mapId uint32, f Operator) {
		ForInMap(l)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
		monsters, err := GetInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(monsters)
	}
}

func GetInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
		resp, err := requestInMap(l)(worldId, channelId, mapId)
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
}

func GetById(l logrus.FieldLogger) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		resp, err := requestById(l)(id)
		if err != nil {
			return nil, err
		}

		d := resp.Data()
		n := makeMonster(id, d.Attributes)
		return &n, nil
	}
}

func makeMonster(id uint32, att attributes) Model {
	return NewMonster(id, att.ControlCharacterId, att.MonsterId, att.X, att.Y, att.Stance, att.FH, att.Team)
}
