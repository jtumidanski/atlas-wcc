package _map

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
		resp, err := requestCharactersInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return nil, err
		}

		cIds := make([]uint32, 0)
		for _, d := range resp.DataList() {
			cId, err := strconv.ParseUint(d.Id, 10, 32)
			if err != nil {
				break
			}
			cIds = append(cIds, uint32(cId))
		}
		return cIds, nil
	}
}