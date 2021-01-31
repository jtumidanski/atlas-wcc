package domain

type NPC struct {
   objectId uint32
   id       uint32
   x        int16
   cy       int16
   f        uint32
   fh       uint16
   rx0      int16
   rx1      int16
}

func NewNPC(objectId uint32, id uint32, x int16, cy int16, f uint32, fh uint16, rx0 int16, rx1 int16) NPC {
   return NPC{
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

func (n NPC) ObjectId() uint32 {
   return n.objectId
}

func (n NPC) Id() uint32 {
   return n.id
}

func (n NPC) X() int16 {
   return n.x
}

func (n NPC) CY() int16 {
   return n.cy
}

func (n NPC) F() uint32 {
   return n.f
}

func (n NPC) Fh() uint16 {
   return n.fh
}

func (n NPC) RX0() int16 {
   return n.rx0
}

func (n NPC) RX1() int16 {
   return n.rx1
}
