package skill

type Model struct {
   id          uint32
   level       uint32
   masterLevel uint32
   expiration  int64
   hidden      bool
   fourthJob   bool
}

func NewSkill(id uint32, level uint32, masterLevel uint32, expiration int64, hidden bool, fourthJob bool) Model {
   return Model{
      id:          id,
      level:       level,
      masterLevel: masterLevel,
      expiration:  expiration,
      hidden:      hidden,
      fourthJob:   fourthJob,
   }
}

func (s Model) Hidden() bool {
   return s.hidden
}

func (s Model) Id() uint32 {
   return s.id
}

func (s Model) Level() uint32 {
   return s.level
}

func (s Model) Expiration() int64 {
   return s.expiration
}

func (s Model) FourthJob() bool {
   return s.fourthJob
}

func (s Model) MasterLevel() uint32 {
   return s.masterLevel
}
