package monster

type Model struct {
   uniqueId           uint32
   controlCharacterId uint32
   monsterId          uint32
   x                  int16
   y                  int16
   stance             byte
   fh                 int16
   team               int8
}

func NewMonster(uniqueId uint32, controlCharacterId uint32, monsterId uint32, x int16, y int16, stance byte, fh int16, team int8) Model {
   return Model{
      uniqueId:           uniqueId,
      controlCharacterId: controlCharacterId,
      monsterId:          monsterId,
      x:                  x,
      y:                  y,
      stance:             stance,
      fh:                 fh,
      team:               team,
   }
}

func (m Model) UniqueId() uint32 {
   return m.uniqueId
}

func (m Model) Controlled() bool {
   return m.controlCharacterId != 0
}

func (m Model) MonsterId() uint32 {
   return m.monsterId
}

func (m Model) X() int16 {
   return m.x
}

func (m Model) Y() int16 {
   return m.y
}

func (m Model) Stance() byte {
   return m.stance
}

func (m Model) FH() int16 {
   return m.fh
}

func (m Model) Team() int8 {
   return m.team
}
