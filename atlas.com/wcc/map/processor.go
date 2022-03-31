package _map

import (
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) []uint32 {
	return func(worldId byte, channelId byte, mapId uint32) []uint32 {
		cIds := make([]uint32, 0)
		resp, err := requestCharactersInMap(worldId, channelId, mapId)(l, span)
		if err != nil {
			return cIds
		}

		for _, d := range resp.DataList() {
			cId, err := strconv.ParseUint(d.Id, 10, 32)
			if err != nil {
				break
			}
			cIds = append(cIds, uint32(cId))
		}
		return cIds
	}
}

func GetOtherCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32) []uint32 {
	return func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32) []uint32 {
		result := make([]uint32, 0)
		ids := GetCharacterIdsInMap(l, span)(worldId, channelId, mapId)
		for _, id := range ids {
			if id != referenceCharacterId {
				result = append(result, id)
			}
		}
		return result
	}
}

func ForSessionsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, o model.Operator[session.Model]) {
	return func(worldId byte, channelId byte, mapId uint32, o model.Operator[session.Model]) {
		session.ForEachByCharacterId(GetCharacterIdsInMap(l, span)(worldId, channelId, mapId), o)
	}
}

func ForOtherSessionsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32, o model.Operator[session.Model]) {
	return func(worldId byte, channelId byte, mapId uint32, referenceCharacterId uint32, o model.Operator[session.Model]) {
		session.ForEachByCharacterId(GetOtherCharacterIdsInMap(l, span)(worldId, channelId, mapId, referenceCharacterId), o)
	}
}
