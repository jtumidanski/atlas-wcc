package domain

type Skill struct {
   id          uint32
   level       uint32
   masterLevel uint32
   expiration  int64
   hidden      bool
   fourthJob   bool
}

func NewSkill(id uint32, level uint32, masterLevel uint32, expiration int64, hidden bool, fourthJob bool) Skill {
   return Skill{
      id:          id,
      level:       level,
      masterLevel: masterLevel,
      expiration:  expiration,
      hidden:      hidden,
      fourthJob:   fourthJob,
   }
}

func (s Skill) Hidden() bool {
   return s.hidden
}

func (s Skill) Id() uint32 {
   return s.id
}

func (s Skill) Level() uint32 {
   return s.level
}

func (s Skill) Expiration() int64 {
   return s.expiration
}

func (s Skill) FourthJob() bool {
   return s.fourthJob
}

func (s Skill) MasterLevel() uint32 {
   return s.masterLevel
}
