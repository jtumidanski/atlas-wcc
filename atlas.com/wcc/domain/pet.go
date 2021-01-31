package domain

type Pet struct {
   id     uint64
   itemId uint32
   name   string
   x      int16
   y      int16
   stance byte
   fh     uint32
}

func (p Pet) Id() uint64 {
   return p.id
}

func (p Pet) ItemId() uint32 {
   return p.itemId
}

func (p Pet) Name() string {
   return p.name
}

func (p Pet) X() int16 {
   return p.x
}

func (p Pet) Y() int16 {
   return p.y
}

func (p Pet) Stance() byte {
   return p.stance
}

func (p Pet) Fh() uint32 {
   return p.fh
}
