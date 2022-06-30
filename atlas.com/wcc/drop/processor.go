package drop

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ForEachInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, f model.Operator[Model]) {
	return func(worldId byte, channelId byte, mapId uint32, f model.Operator[Model]) {
		model.ForEach(InMapModelProvider(l, span)(worldId, channelId, mapId), f)
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[Model] {
	return func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(worldId, channelId, mapId), makeModel)
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.Atoi(body.Id)
	if err != nil {
		return Model{}, err
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
	return m, nil
}

func SpawnDropForSession(l logrus.FieldLogger) func(s session.Model) model.Operator[Model] {
	return func(s session.Model) model.Operator[Model] {
		return func(d Model) error {
			var a = uint32(0)
			if d.ItemId() != 0 {
				a = 0
			} else {
				a = d.Meso()
			}
			err := session.AnnounceOperator(WriteDropItemFromMapObject(l)(d.UniqueId(), d.ItemId(), d.Meso(), a,
				d.DropperUniqueId(), d.DropType(), d.OwnerId(), d.OwnerPartyId(), s.CharacterId(),
				0, d.DropTime(), d.DropX(), d.DropY(), d.DropperX(), d.DropperY(),
				d.CharacterDrop(), d.Mod()))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce drop to character %d", s.CharacterId())
			}
			return err
		}
	}
}
