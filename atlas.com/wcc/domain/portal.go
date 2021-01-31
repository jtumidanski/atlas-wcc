package domain

type Portal struct {
   id          uint32
   mapId       uint32
   name        string
   target      string
   targetMapId uint32
   theType     uint32
   x           int32
   y           int32
   scriptName  string
}

func (p Portal) Id() uint32 {
   return p.id
}

func NewPortal(id uint32, mapId uint32, name string, target string, targetMapId uint32, theType uint32, x int32, y int32, scriptName string) Portal {
   return Portal{
      id:          id,
      mapId:       mapId,
      name:        name,
      target:      target,
      targetMapId: targetMapId,
      theType:     theType,
      x:           x,
      y:           y,
      scriptName:  scriptName,
   }
}
