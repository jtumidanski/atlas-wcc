package reactor

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ForEachAliveInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f model.Operator[Model]) {
	return func(worldId byte, channelId byte, mapId uint32, f model.Operator[Model]) {
		model.ForEach(InMapModelProvider(l, span)(worldId, channelId, mapId, AliveFilter), f)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...model.Filter[Model]) model.SliceProvider[Model] {
	return func(worldId byte, channelId byte, mapId uint32, filters ...model.Filter[Model]) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(worldId, channelId, mapId), makeModel, filters...)
	}
}

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(id uint32) model.Provider[Model] {
	return func(id uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(id), makeModel)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(id uint32) (Model, error) {
	return func(id uint32) (Model, error) {
		return ByIdModelProvider(l, span)(id)()
	}
}

func AliveFilter(m Model) bool {
	return m.Alive()
}

func makeModel(data requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(data.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	attr := data.Attributes
	return Model{
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

func SpawnForSession(l logrus.FieldLogger) func(s session.Model) model.Operator[Model] {
	return func(s session.Model) model.Operator[Model] {
		return func(r Model) error {
			err := session.Announce(WriteReactorSpawn(l)(r.Id(), r.Classification(), r.State(), r.X(), r.Y()))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d creation to session %d.", r.Id(), s.SessionId())
			}
			return err
		}
	}
}

func DestroyForSession(l logrus.FieldLogger, span opentracing.Span) func(reactorId uint32) model.Operator[session.Model] {
	return func(reactorId uint32) model.Operator[session.Model] {
		r, err := GetById(l, span)(reactorId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate reactor to process status of.")
			return model.ErrorOperator[session.Model](err)
		}
		return session.Announce(WriteReactorDestroyed(l)(r.Id(), r.State(), r.X(), r.Y()))
	}
}

func HitForSession(l logrus.FieldLogger, span opentracing.Span) func(reactorId uint32, stance uint16) model.Operator[session.Model] {
	return func(reactorId uint32, stance uint16) model.Operator[session.Model] {
		r, err := GetById(l, span)(reactorId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate reactor to process status of.")
			return model.ErrorOperator[session.Model](err)
		}
		return session.Announce(WriteReactorTrigger(l)(r.Id(), r.State(), r.X(), r.Y(), byte(stance)))
	}
}

func CreateForSession(l logrus.FieldLogger, span opentracing.Span) func(reactorId uint32, stance uint16) model.Operator[session.Model] {
	return func(reactorId uint32, stance uint16) model.Operator[session.Model] {
		r, err := GetById(l, span)(reactorId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate reactor to process status of.")
			return model.ErrorOperator[session.Model](err)
		}
		return session.Announce(WriteReactorSpawn(l)(r.Id(), r.Classification(), r.State(), r.X(), r.Y()))
	}
}
