package consumers

import (
   "atlas-wcc/kafka/handler"
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "atlas-wcc/socket/response/writer"
   "github.com/sirupsen/logrus"
)

type rangeAttackEvent struct {
   WorldId            byte                `json:"worldId"`
   ChannelId          byte                `json:"channelId"`
   MapId              uint32              `json:"mapId"`
   CharacterId        uint32              `json:"characterId"`
   SkillId            uint32              `json:"skillId"`
   SkillLevel         byte                `json:"skillLevel"`
   Stance             byte                `json:"stance"`
   AttackedAndDamaged byte                `json:"attackedAndDamaged"`
   Projectile         uint32              `json:"projectile"`
   Damage             map[uint32][]uint32 `json:"damage"`
   Speed              byte                `json:"speed"`
   Direction          byte                `json:"direction"`
   Display            byte                `json:"display"`
}

func EmptyRangeAttackEventCreator() handler.EmptyEventCreator {
   return func() interface{} {
      return &rangeAttackEvent{}
   }
}

func HandleRangeAttackEvent() ChannelEventProcessor {
   return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
      if event, ok := e.(*rangeAttackEvent); ok {
         if wid != event.WorldId || cid != event.ChannelId {
            return
         }

         processors.ForEachSessionInMap(event.WorldId, event.ChannelId, event.MapId, writeRangeAttack(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display, event.Projectile))
      } else {
         l.Errorf("Unable to cast event provided to handler")
      }
   }
}

func writeRangeAttack(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, projectile uint32) func(s mapleSession.MapleSession) {
   return func(s mapleSession.MapleSession) {
      s.Announce(writer.WriteRangeAttack(characterId, skill, skillLevel, stance, numberAttackedAndDamaged, damage, speed, direction, display, projectile))
   }
}
