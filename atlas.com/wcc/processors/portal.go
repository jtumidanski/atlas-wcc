package processors

import (
   "atlas-wcc/domain"
   "atlas-wcc/rest/attributes"
   "atlas-wcc/rest/requests"
   "strconv"
)

func GetPortalByName(mapId uint32, portalName string) (*domain.Portal, error) {
   resp, err := requests.MapInformation().GetPortalByName(mapId, portalName)
   if err != nil {
      return nil, err
   }

   d := resp.Data()
   aid, err := strconv.ParseUint(d.Id, 10, 32)
   if err != nil {
      return nil, err
   }

   a := makePortal(uint32(aid), mapId, d.Attributes)
   return &a, nil
}

func makePortal(id uint32, mapId uint32, attr attributes.PortalAttributes) domain.Portal {
   return domain.NewPortal(id, mapId, attr.Name, attr.Target, attr.TargetMap, attr.Type, attr.X, attr.Y, attr.ScriptName)
}
