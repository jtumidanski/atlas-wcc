package reactor

import (
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelOperator func(*Model)

type ModelListOperator func([]*Model)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

func requestModelProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes]) ModelProvider {
	return func(r requests.Request[attributes]) ModelProvider {
		return func() (*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			p, err := makeModel(resp.Data())
			if err != nil {
				return nil, err
			}
			return p, nil
		}
	}
}

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes], filters ...Filter) ModelListProvider {
	return func(r requests.Request[attributes], filters ...Filter) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(v)
				if err != nil {
					return nil, err
				}
				ok := true
				for _, filter := range filters {
					if !filter(m) {
						ok = false
						break
					}
				}
				if ok {
					ms = append(ms, m)
				}
			}
			return ms, nil
		}
	}
}

func ExecuteForEach(f ModelOperator) ModelListOperator {
	return func(drops []*Model) {
		for _, drop := range drops {
			f(drop)
		}
	}
}

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f ModelOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f ModelOperator) {
		ForReactorsInMap(l, span)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForEachAliveInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f ModelOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f ModelOperator) {
		ForAliveReactorsInMap(l, span)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForReactorsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f ModelListOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f ModelListOperator) {
		reactors, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(reactors)
	}
}

func ForAliveReactorsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f ModelListOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f ModelListOperator) {
		reactors, err := GetInMap(l, span)(worldId, channelId, mapId, AliveFilter())
		if err != nil {
			return
		}
		f(reactors)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ModelListProvider {
	return func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ModelListProvider {
		return requestModelListProvider(l, span)(requestInMap(worldId, channelId, mapId), filters...)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ([]*Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ([]*Model, error) {
		return InMapModelProvider(l, span)(worldId, channelId, mapId, filters...)()
	}
}

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(id uint32) ModelProvider {
	return func(id uint32) ModelProvider {
		return requestModelProvider(l, span)(requestById(id))
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		return ByIdModelProvider(l, span)(id)()
	}
}

type Filter func(*Model) bool

func AliveFilter() Filter {
	return func(m *Model) bool {
		return m.Alive()
	}
}

func makeModel(data requests.DataBody[attributes]) (*Model, error) {
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

func SpawnForSession(l logrus.FieldLogger) func(s *session.Model) ModelOperator {
	return func(s *session.Model) ModelOperator {
		return func(r *Model) {
			err := s.Announce(WriteReactorSpawn(l)(r.Id(), r.Classification(), r.State(), r.X(), r.Y()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d creation to session %d.", r.Id(), s.SessionId())
			}
		}
	}
}

func DestroyForSession(l logrus.FieldLogger, span opentracing.Span) func(reactorId uint32) session.Operator {
	return func(reactorId uint32) session.Operator {
		return func(session *session.Model) {
			r, err := GetById(l, span)(reactorId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate reactor to process status of.")
				return
			}
			err = session.Announce(WriteReactorDestroyed(l)(r.Id(), r.State(), r.X(), r.Y()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d destroyed to session %d.", r.Id(), session.SessionId())
			}
		}
	}
}

func HitForSession(l logrus.FieldLogger, span opentracing.Span) func(reactorId uint32, stance uint16) session.Operator {
	return func(reactorId uint32, stance uint16) session.Operator {
		return func(session *session.Model) {
			r, err := GetById(l, span)(reactorId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate reactor to process status of.")
				return
			}
			err = session.Announce(WriteReactorTrigger(l)(r.Id(), r.State(), r.X(), r.Y(), byte(stance)))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d trigger to session %d.", r.Id(), session.SessionId())
			}
		}
	}
}

func CreateForSession(l logrus.FieldLogger, span opentracing.Span) func(reactorId uint32, stance uint16) session.Operator {
	return func(reactorId uint32, stance uint16) session.Operator {
		return func(session *session.Model) {
			r, err := GetById(l, span)(reactorId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate reactor to process status of.")
				return
			}
			err = session.Announce(WriteReactorSpawn(l)(r.Id(), r.Classification(), r.State(), r.X(), r.Y()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d creation to session %d.", r.Id(), session.SessionId())
			}
		}
	}
}
