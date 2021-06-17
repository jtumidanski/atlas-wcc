package npc

type Model struct {
   objectId uint32
   id       uint32
   x        int16
   cy       int16
   f        uint32
   fh       uint16
   rx0      int16
   rx1      int16
}

func NewNPC(objectId uint32, id uint32, x int16, cy int16, f uint32, fh uint16, rx0 int16, rx1 int16) Model {
   return Model{
      objectId: objectId,
      id:       id,
      x:        x,
      cy:       cy,
      f:        f,
      fh:       fh,
      rx0:      rx0,
      rx1:      rx1,
   }
}

func (n Model) ObjectId() uint32 {
   return n.objectId
}

func (n Model) Id() uint32 {
   return n.id
}

func (n Model) X() int16 {
   return n.x
}

func (n Model) CY() int16 {
   return n.cy
}

func (n Model) F() uint32 {
   return n.f
}

func (n Model) Fh() uint16 {
   return n.fh
}

func (n Model) RX0() int16 {
   return n.rx0
}

func (n Model) RX1() int16 {
   return n.rx1
}
