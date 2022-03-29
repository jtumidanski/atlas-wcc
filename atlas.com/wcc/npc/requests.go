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

func requestNPCsInMap(mapId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(npcsInMap, mapId))
}

func requestNPCsInMapByObjectId(mapId uint32, objectId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(npcsInMapByObjectId, mapId, objectId))
}
