package reactor

import (
	"github.com/sirupsen/logrus"
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

func ForEachInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, f Operator) {
	return func(worldId byte, channelId byte, mapId uint32, f Operator) {
		ForReactorsInMap(l)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForReactorsInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
		reactors, err := GetInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(reactors)
	}
}

func GetById(l logrus.FieldLogger) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		data, err := requestById(l)(id)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve reactor by id %d.", id)
			return nil, err
		}

		r, err := makeReactor(data.Data)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}

func GetInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
		resp, err := requestInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return nil, err
		}

		reactors := make([]Model, 0)
		for _, d := range resp.Data {
			r, err := makeReactor(d)
			if err != nil {
				l.WithError(err).Errorf("Unable to make reactor %d model.", d.Attributes.Classification)
			} else {
				reactors = append(reactors, *r)
			}
		}
		return reactors, nil
	}
}

func makeReactor(data DataBody) (*Model, error) {
	id, err := strconv.ParseUint(data.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	attr := data.Attributes
	return &Model{
		id:             uint32(id),
		classification: attr.Classification,
		name:           attr.Name,
		state:          attr.State,
		eventState:     attr.EventState,
		delay:          attr.Delay,
		direction:      attr.FacingDirection,
		x:              attr.X,
		y:              attr.Y,
	}, nil
}
