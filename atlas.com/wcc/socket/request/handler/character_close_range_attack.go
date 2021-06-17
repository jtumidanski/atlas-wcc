package handler

import (
   "atlas-wcc/character"
   "atlas-wcc/kafka/producers"
   "atlas-wcc/session"
   request2 "atlas-wcc/socket/request"
   "github.com/jtumidanski/atlas-socket/request"
   "github.com/sirupsen/logrus"
   "math"
)

const OpCharacterCloseRangeAttack uint16 = 0x2C

type attackPacket struct {
   numberAttacked           byte
   numberDamaged            byte
   numberAttackedAndDamaged byte
   skill                    uint32
   skillLevel               byte
   stance                   byte
   direction                byte
   rangedDirection          byte
   charge                   uint32
   display                  byte
   ranged                   bool
   magic                    bool
   speed                    byte
   allDamage                map[uint32][]uint32
   x                        int16
   y                        int16
}

func (p attackPacket) Skill() uint32 {
   return p.skill
}

func (p attackPacket) SkillLevel() byte {
   return p.skillLevel
}

func (p attackPacket) Stance() byte {
   return p.stance
}

func (p attackPacket) NumberAttackedAndDamaged() byte {
   return p.numberAttackedAndDamaged
}

func (p attackPacket) AllDamage() map[uint32][]uint32 {
   return p.allDamage
}

func (p attackPacket) Speed() byte {
   return p.speed
}

func (p attackPacket) Direction() byte {
   return p.direction
}

func (p attackPacket) Display() byte {
   return p.display
}

func (p attackPacket) RangedDirection() byte {
   return p.rangedDirection
}

func (p attackPacket) NumberAttacked() byte {
   return p.numberAttacked
}

func (p attackPacket) NumberDamaged() byte {
   return p.numberDamaged
}

func (p attackPacket) Charge() uint32 {
   return p.charge
}

func (p attackPacket) Ranged() bool {
   return p.ranged
}

func (p attackPacket) Magic() bool {
   return p.magic
}

func (p attackPacket) X() int16 {
   return p.x
}

func (p attackPacket) Y() int16 {
   return p.y
}

func readAttackPacket(reader *request.RequestReader, characterId uint32, ranged bool, magic bool) attackPacket {
   reader.ReadByte()
   numberAttackedAndDamaged := reader.ReadByte()
   numberAttacked := numberAttackedAndDamaged >> 4 & 0xF
   numberDamaged := numberAttackedAndDamaged & 0xF
   skillId := reader.ReadUint32()
   skillLevel := byte(0)
   charge := uint32(0)
   reader.Skip(8)
   display := reader.ReadByte()
   direction := reader.ReadByte()
   stance := reader.ReadByte()
   var speed byte
   var rangedDirection byte
   if ranged {
      reader.ReadByte()
      speed = reader.ReadByte()
      reader.ReadByte()
      rangedDirection = reader.ReadByte()
      reader.Skip(7)
   } else {
      reader.ReadByte()
      speed = reader.ReadByte()
      reader.Skip(4)
   }
   calculatedMaximumDamage := character.GetCharacterWeaponDamage(characterId)
   bonusDamageBuff := uint32(100)
   if bonusDamageBuff != 100 {
      damageBuff := bonusDamageBuff / 100
      calculatedMaximumDamage = uint32(math.Float64bits(math.Ceil(float64(calculatedMaximumDamage * damageBuff))))
   }
   canCritical := false
   shadowPartner := false

   damage := make(map[uint32][]uint32, 0)
   for i := byte(0); i < numberAttacked; i++ {
      oid := reader.ReadUint32()
      reader.Skip(14)

      var allDamageNumbers []uint32
      for j := byte(0); j < numberDamaged; j++ {
         damage := reader.ReadUint32()
         hitDamageMax := calculatedMaximumDamage
         if shadowPartner {
            if j >= numberDamaged/2 {
               hitDamageMax = uint32(math.Ceil(float64(hitDamageMax) * 0.5))
            }
         }
         maxWithCritical := hitDamageMax
         if canCritical {
            maxWithCritical *= 2
         }
         allDamageNumbers = append(allDamageNumbers, damage)
      }
      reader.Skip(4)
      damage[oid] = allDamageNumbers
   }
   return attackPacket{
      numberAttacked:           numberAttacked,
      numberDamaged:            numberDamaged,
      numberAttackedAndDamaged: numberAttackedAndDamaged,
      skill:                    skillId,
      skillLevel:               skillLevel,
      stance:                   stance,
      direction:                direction,
      rangedDirection:          rangedDirection,
      charge:                   charge,
      display:                  display,
      ranged:                   ranged,
      magic:                    magic,
      speed:                    speed,
      allDamage:                damage,
      x:                        0,
      y:                        0,
   }
}

func CharacterCloseRangeAttackHandler() request2.MessageHandler {
   return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
      p := readAttackPacket(r, s.CharacterId(), false, false)

      catt, err := character.GetCharacterAttributesById(s.CharacterId())
      if err != nil {
         l.WithError(err).Errorf("Unable to retrieve character attributes for character %d.", s.CharacterId())
         return
      }
      producers.CharacterAttack(l)(s.WorldId(), s.ChannelId(), catt.MapId(), s.CharacterId(), p.Skill(), p.SkillLevel(), p.NumberAttacked(), p.NumberDamaged(), p.NumberAttackedAndDamaged(), p.Stance(), p.Direction(), p.RangedDirection(), p.Charge(), p.Display(), p.Ranged(), p.Magic(), p.Speed(), p.AllDamage(), p.X(), p.Y())
   }
}
