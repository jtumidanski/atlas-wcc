package monster

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

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes]) ModelListProvider {
	return func(r requests.Request[attributes]) ModelListProvider {
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
				ms = append(ms, m)
			}
			return ms, nil
		}
	}
}

func ExecuteForEach(f ModelOperator) ModelListOperator {
	return func(monsters []*Model) {
		for _, monster := range monsters {
			f(monster)
		}
	}
}

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f ModelOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f ModelOperator) {
		ForInMap(l, span)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f ModelListOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f ModelListOperator) {
		monsters, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(monsters)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ModelListProvider {
	return func(worldId byte, channelId byte, mapId uint32) ModelListProvider {
		return requestModelListProvider(l, span)(requestInMap(worldId, channelId, mapId))
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ([]*Model, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]*Model, error) {
		return InMapModelProvider(l, span)(worldId, channelId, mapId)()
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

func makeModel(body requests.DataBody[attributes]) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := NewMonster(uint32(id), att.ControlCharacterId, att.MonsterId, att.X, att.Y, att.Stance, att.FH, att.Team)
	return &m, nil
}

func SpawnForSession(l logrus.FieldLogger) func(s *session.Model) ModelOperator {
	return func(s *session.Model) ModelOperator {
		return func(m *Model) {
			err := s.Announce(WriteSpawnMonster(l)(m, false))
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn monster %d for character %d", m.MonsterId(), s.CharacterId())
			}
		}
	}
}

func DestroyForSession(l logrus.FieldLogger, uniqueId uint32) session.Operator {
	k1 := WriteKillMonster(l)(uniqueId, false)
	k2 := WriteKillMonster(l)(uniqueId, true)
	return func(s *session.Model) {
		err := s.Announce(k1)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = s.Announce(k2)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func CreateForSession(l logrus.FieldLogger, m *Model) session.Operator {
	sm := WriteSpawnMonster(l)(m, false)
	return func(s *session.Model) {
		err := s.Announce(sm)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func KillForSession(l logrus.FieldLogger, uniqueId uint32) session.Operator {
	b := WriteKillMonster(l)(uniqueId, true)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func MoveForSession(l logrus.FieldLogger, objectId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, movementList []byte) session.Operator {
	b := WriteMoveMonster(l)(objectId, skillPossible, skill, skillId, skillLevel, option, startX, startY, movementList)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
