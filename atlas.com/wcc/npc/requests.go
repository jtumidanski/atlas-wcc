package npc

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	npcsInMap                          = mapsResource + "%d/npcs"
	npcsInMapByObjectId                = mapsResource + "%d/npcs?objectId=%d"
)

func requestNPCsInMap(mapId uint32) (*dataContainer, error) {
	ar := &dataContainer{}
	err := requests.Get(fmt.Sprintf(npcsInMap, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func requestNPCsInMapByObjectId(mapId uint32, objectId uint32) (*dataContainer, error) {
	ar := &dataContainer{}
	err := requests.Get(fmt.Sprintf(npcsInMapByObjectId, mapId, objectId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
