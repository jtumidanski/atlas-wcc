package pet

type Model struct {
   id     uint64
   itemId uint32
   name   string
   x      int16
   y      int16
   stance byte
   fh     uint32
}

func (p Model) Id() uint64 {
   return p.id
}

func (p Model) ItemId() uint32 {
   return p.itemId
}

func (p Model) Name() string {
   return p.name
}

func (p Model) X() int16 {
   return p.x
}

func (p Model) Y() int16 {
   return p.y
}

func (p Model) Stance() byte {
   return p.stance
}

func (p Model) Fh() uint32 {
   return p.fh
}
