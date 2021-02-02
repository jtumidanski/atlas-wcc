package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = baseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	npcsInMap                          = mapsResource + "%d/npcs"
	npcsInMapByObjectId                = mapsResource + "%d/npcs?objectId=%d"
	portalsResource                    = mapsResource + "%d/portals"
	portalsByName                      = portalsResource + "?name=%s"
)

var MapInformation = func() *mapInformation {
	return &mapInformation{}
}

type mapInformation struct {
}

func (m *mapInformation) GetPortalByName(mapId uint32, portalName string) (*attributes.PortalDataContainer, error) {
	ar := &attributes.PortalDataContainer{}
	err := get(fmt.Sprintf(portalsByName, mapId, portalName), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (m *mapInformation) GetNPCsInMap(mapId uint32) (*attributes.NpcDataContainer, error) {
	ar := &attributes.NpcDataContainer{}
	err := get(fmt.Sprintf(npcsInMap, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (m *mapInformation) GetNPCsInMapByObjectId(mapId uint32, objectId uint32) (*attributes.NpcDataContainer, error) {
	ar := &attributes.NpcDataContainer{}
	err := get(fmt.Sprintf(npcsInMapByObjectId, mapId, objectId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
