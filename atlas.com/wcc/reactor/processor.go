package reactor

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

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
