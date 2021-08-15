package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/reactor"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type reactorStatusEvent struct {
	WorldId   byte   `json:"world_id"`
	ChannelId byte   `json:"channel_id"`
	MapId     uint32 `json:"map_id"`
	Id        uint32 `json:"id"`
	Status    string `json:"status"`
	Stance    uint16 `json:"stance"`
}

func EmptyReactorStatusEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &reactorStatusEvent{}
	}
}

func HandleReactorStatusEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*reactorStatusEvent); ok {
			if event.WorldId != wid || event.ChannelId != cid {
				return
			}
			if event.Status == "CREATED" {
				session.ForEachInMap(l)(wid, cid, event.MapId, handleReactorCreation(l)(event.Id, event.Stance))
			} else if event.Status == "TRIGGERED" {
				session.ForEachInMap(l)(wid, cid, event.MapId, handleReactorHit(l)(event.Id, event.Stance))
			} else if event.Status == "DESTROYED" {
				session.ForEachInMap(l)(wid, cid, event.MapId, handleReactorDestroy(l)(event.Id))
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func handleReactorDestroy(l logrus.FieldLogger) func(reactorId uint32) session.Operator {
	return func(reactorId uint32) session.Operator {
		return func(session *session.Model) {
			r, err := reactor.GetById(l)(reactorId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate reactor to process status of.")
				return
			}
			err = session.Announce(writer.WriteReactorDestroyed(l)(r.Id(), r.State(), r.X(), r.Y()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d destroyed to session %d.", r.Id(), session.SessionId())
			}
		}
	}
}

func handleReactorHit(l logrus.FieldLogger) func(reactorId uint32, stance uint16) session.Operator {
	return func(reactorId uint32, stance uint16) session.Operator {
		return func(session *session.Model) {
			r, err := reactor.GetById(l)(reactorId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate reactor to process status of.")
				return
			}
			err = session.Announce(writer.WriteReactorTrigger(l)(r.Id(), r.State(), r.X(), r.Y(), byte(stance)))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d trigger to session %d.", r.Id(), session.SessionId())
			}
		}
	}
}

func handleReactorCreation(l logrus.FieldLogger) func(reactorId uint32, stance uint16) session.Operator {
	return func(reactorId uint32, stance uint16) session.Operator {
		return func(session *session.Model) {
			r, err := reactor.GetById(l)(reactorId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate reactor to process status of.")
				return
			}
			err = session.Announce(writer.WriteReactorSpawn(l)(r.Id(), r.Classification(), r.State(), r.X(), r.Y()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d creation to session %d.", r.Id(), session.SessionId())
			}
		}
	}
}
