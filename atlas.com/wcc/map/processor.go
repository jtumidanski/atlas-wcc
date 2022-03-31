package _map

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func CharacterIdsInMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[uint32] {
	return func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[uint32] {
		return requests.SliceProvider[characterAttributes, uint32](l, span)(requestCharactersInMap(worldId, channelId, mapId), getCharacterId)
	}
}

func GetCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
		return CharacterIdsInMapModelProvider(l, span)(worldId, channelId, mapId)()
	}
}

func getCharacterId(body requests.DataBody[characterAttributes]) (uint32, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func NotCharacterIdFilter(referenceCharacterId uint32) func(characterId uint32) bool {
	return func(characterId uint32) bool {
		return referenceCharacterId != characterId
	}
}

func OtherCharacterIdsInMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32) model.SliceProvider[uint32] {
	return func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32) model.SliceProvider[uint32] {
		return model.FilteredProvider(CharacterIdsInMapModelProvider(l, span)(worldId, channelId, mapId), NotCharacterIdFilter(referenceCharacterId))
	}
}

func GetOtherCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32) ([]uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32) ([]uint32, error) {
		return OtherCharacterIdsInMapModelProvider(l, span)(worldId, channelId, mapId, referenceCharacterId)()
	}
}

func ForSessionsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, o model.Operator[session.Model]) {
	return func(worldId byte, channelId byte, mapId uint32, o model.Operator[session.Model]) {
		ids, err := GetCharacterIdsInMap(l, span)(worldId, channelId, mapId)
		if err == nil {
			session.ForEachByCharacterId(ids, o)
		}
	}
}

func ForOtherSessionsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32, o model.Operator[session.Model]) {
	return func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32, o model.Operator[session.Model]) {
		ids, err := GetOtherCharacterIdsInMap(l, span)(worldId, channelId, mapId, referenceCharacterId)
		if err == nil {
			session.ForEachByCharacterId(ids, o)
		}
	}
}
