package portal

type Model struct {
   id          uint32
   mapId       uint32
   name        string
   target      string
   targetMapId uint32
   theType     uint8
   x           int16
   y           int16
   scriptName  string
}

func (p Model) Id() uint32 {
   return p.id
}

func NewPortal(id uint32, mapId uint32, name string, target string, targetMapId uint32, theType uint8, x int16, y int16, scriptName string) Model {
   return Model{
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
