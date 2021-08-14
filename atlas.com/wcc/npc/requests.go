package npc

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	npcsInMap                          = mapsResource + "%d/npcs"
	npcsInMapByObjectId                = mapsResource + "%d/npcs?objectId=%d"
)

func requestNPCsInMap(l logrus.FieldLogger) func(mapId uint32) (*dataContainer, error) {
	return func(mapId uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(npcsInMap, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestNPCsInMapByObjectId(l logrus.FieldLogger) func(mapId uint32, objectId uint32) (*dataContainer, error) {
	return func(mapId uint32, objectId uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(npcsInMapByObjectId, mapId, objectId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
