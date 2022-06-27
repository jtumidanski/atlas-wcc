package _map

import (
	"atlas-wcc/model"
	"atlas-wcc/portal"
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

func ForSessionsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, o model.Operator[session.Model]) {
	return func(worldId byte, channelId byte, mapId uint32, o model.Operator[session.Model]) {
		session.ForEachByCharacterId(CharacterIdsInMapModelProvider(l, span)(worldId, channelId, mapId), o)
	}
}

func ForOtherSessionsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32, o model.Operator[session.Model]) {
	return func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32, o model.Operator[session.Model]) {
		session.ForEachByCharacterId(OtherCharacterIdsInMapModelProvider(l, span)(worldId, channelId, mapId, referenceCharacterId), o)
	}
}

func WarpRandom(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, mapId uint32) {
	return func(worldId byte, channelId byte, characterId uint32, mapId uint32) {
		WarpToPortal(l, span)(worldId, channelId, characterId, mapId, portal.RandomPortalIdProvider(l, span)(mapId))
	}
}

func WarpToPortal(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, mapId uint32, p portal.IdProvider) {
	return func(worldId byte, channelId byte, characterId uint32, mapId uint32, p portal.IdProvider) {
		emitChangeMap(l, span)(worldId, channelId, characterId, mapId, p())
	}
}
