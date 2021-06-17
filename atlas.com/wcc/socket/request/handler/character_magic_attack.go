package handler

import (
   "atlas-wcc/character"
   "atlas-wcc/kafka/producers"
   "atlas-wcc/session"
   request2 "atlas-wcc/socket/request"
   "github.com/jtumidanski/atlas-socket/request"
   "github.com/sirupsen/logrus"
)

const OpCharacterMagicAttack uint16 = 0x2E

func CharacterMagicAttackHandler() request2.MessageHandler {
   return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
      p := readAttackPacket(r, s.CharacterId(), false, true)

      catt, err := character.GetCharacterAttributesById(s.CharacterId())
      if err != nil {
         l.WithError(err).Errorf("Unable to retrieve character attributes for character %d.", s.CharacterId())
         return
      }
      producers.CharacterAttack(l)(s.WorldId(), s.ChannelId(), catt.MapId(), s.CharacterId(), p.Skill(), p.SkillLevel(), p.NumberAttacked(), p.NumberDamaged(), p.NumberAttackedAndDamaged(), p.Stance(), p.Direction(), p.RangedDirection(), p.Charge(), p.Display(), p.Ranged(), p.Magic(), p.Speed(), p.AllDamage(), p.X(), p.Y())
   }
}