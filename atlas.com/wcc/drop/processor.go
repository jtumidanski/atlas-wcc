package drop

import (
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Operator func(*Model)

type SliceOperator func([]*Model)

func ExecuteForEach(f Operator) SliceOperator {
	return func(drops []*Model) {
		for _, drop := range drops {
			f(drop)
		}
	}
}

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f Operator) {
	return func(worldId byte, channelId byte, mapId uint32, f Operator) {
		ForDropsInMap(l, span)(worldId, channelId, mapId, ExecuteForEach(f))
	}
}

func ForDropsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
	return func(worldId byte, channelId byte, mapId uint32, f SliceOperator) {
		drops, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return
		}
		f(drops)
	}
}

type ModelListProvider func() ([]*Model, error)

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

func makeModel(body requests.DataBody[attributes]) (*Model, error) {
	id, err := strconv.Atoi(body.Id)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := NewBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetMapId(att.MapId).
		SetUniqueId(uint32(id)).
		SetItemId(att.ItemId).
		SetMeso(att.Meso).
		SetDropType(att.DropType).
		SetDropX(att.DropX).
		SetDropY(att.DropY).
		SetOwnerId(att.OwnerId).
		SetOwnerPartyId(att.OwnerPartyId).
		SetDropTime(att.DropTime).
		SetDropperUniqueId(att.DropperUniqueId).
		SetDropperX(att.DropperX).
		SetDropperY(att.DropperY).
		SetCharacterDrop(att.CharacterDrop).
		SetMod(att.Mod).
		Build()
	return &m, nil
}

func SpawnDropForSession(l logrus.FieldLogger) func(s *session.Model) Operator {
	return func(s *session.Model) Operator {
		return func(d *Model) {
			var a = uint32(0)
			if d.ItemId() != 0 {
				a = 0
			} else {
				a = d.Meso()
			}
			err := s.Announce(WriteDropItemFromMapObject(l)(d.UniqueId(), d.ItemId(), d.Meso(), a,
				d.DropperUniqueId(), d.DropType(), d.OwnerId(), d.OwnerPartyId(), s.CharacterId(),
				0, d.DropTime(), d.DropX(), d.DropY(), d.DropperX(), d.DropperY(),
				d.CharacterDrop(), d.Mod()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce drop to character %d", s.CharacterId())
			}
		}
	}
}

func RemoveDropForSession(l logrus.FieldLogger) func(dropId uint32, characterId uint32) session.Operator {
	return func(dropId uint32, characterId uint32) session.Operator {
		return func(s *session.Model) {
			err := s.Announce(WriteRemoveItem(l)(dropId, 2, characterId))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		}
	}
}

func ExpireDropForSession(l logrus.FieldLogger) func(dropId uint32) session.Operator {
	return func(dropId uint32) session.Operator {
		return func(s *session.Model) {
			err := s.Announce(WriteRemoveItem(l)(dropId, 0, 0))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		}
	}
}
