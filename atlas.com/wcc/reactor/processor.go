package reactor

import (
	"github.com/opentracing/opentracing-go"
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

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f Operator) {
	return func(worldId byte, channelId byte, mapId uint32, f Operator) {
		ForReactorsInMap(l, span)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForEachAliveInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f Operator) {
	return func(worldId byte, channelId byte, mapId uint32, f Operator) {
		ForAliveReactorsInMap(l, span)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForReactorsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
		reactors, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(reactors)
	}
}

func ForAliveReactorsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
		reactors, err := GetInMap(l, span)(worldId, channelId, mapId, AliveFilter())
		if err != nil {
			return
		}
		f(reactors)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		data, err := requestById(l, span)(id)
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

type Filter func(*Model) bool

func AliveFilter() Filter {
	return func(m *Model) bool {
		return m.Alive()
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ([]Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ([]Model, error) {
		resp, err := requestInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return nil, err
		}

		reactors := make([]Model, 0)
		for _, d := range resp.Data {
			r, err := makeReactor(d)
			if err != nil {
				l.WithError(err).Errorf("Unable to make reactor %d model.", d.Attributes.Classification)
			} else {
				ok := true
				for _, filter := range filters {
					if !filter(r) {
						ok = false
						break
					}
				}
				if ok {
					reactors = append(reactors, *r)
				}
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
		alive:          attr.Alive,
	}, nil
}
