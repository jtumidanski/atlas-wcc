package domain

type Monster struct {
   uniqueId           uint32
   controlCharacterId uint32
   monsterId          uint32
   x                  int16
   y                  int16
   stance             byte
   fh                 int16
   team               int8
}

func NewMonster(uniqueId uint32, controlCharacterId uint32, monsterId uint32, x int16, y int16, stance byte, fh int16, team int8) Monster {
   return Monster{
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

func (m Monster) UniqueId() uint32 {
   return m.uniqueId
}

func (m Monster) Controlled() bool {
   return m.controlCharacterId != 0
}

func (m Monster) MonsterId() uint32 {
   return m.monsterId
}

func (m Monster) X() int16 {
   return m.x
}

func (m Monster) Y() int16 {
   return m.y
}

func (m Monster) Stance() byte {
   return m.stance
}

func (m Monster) FH() int16 {
   return m.fh
}

func (m Monster) Team() int8 {
   return m.team
}
