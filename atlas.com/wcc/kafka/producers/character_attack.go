package producers

import "github.com/sirupsen/logrus"

type attackCommand struct {
   WorldId                  byte                `json:"worldId"`
   ChannelId                byte                `json:"channelId"`
   MapId                    uint32              `json:"mapId"`
   CharacterId              uint32              `json:"characterId"`
   NumberAttacked           byte                `json:"numberAttacked"`
   NumberDamaged            byte                `json:"numberDamaged"`
   NumberAttackedAndDamaged byte                `json:"NumberAttackedAndDamaged"`
   SkillId                  uint32              `json:"skillId"`
   SkillLevel               byte                `json:"skillLevel"`
   Stance                   byte                `json:"stance"`
   Direction                byte                `json:"direction"`
   RangedDirection          byte                `json:"rangedDirection"`
   Charge                   uint32              `json:"charge"`
   Display                  byte                `json:"display"`
   Ranged                   bool                `json:"ranged"`
   Magic                    bool                `json:"magic"`
   Speed                    byte                `json:"speed"`
   AllDamage                map[uint32][]uint32 `json:"allDamage"`
   X                        int16               `json:"x"`
   Y                        int16               `json:"y"`
}

func CharacterAttack(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32, skillId uint32, skillLevel byte, attacked byte, damaged byte, attackedAndDamaged byte, stance byte, direction byte, rangedDirection byte, charge uint32, display byte, ranged bool, magic bool, speed byte, allDamage map[uint32][]uint32, x int16, y int16) {
   producer := ProduceEvent(l, "TOPIC_CHARACTER_ATTACK_COMMAND")
   return func(worldId byte, channelId byte, mapId uint32, characterId uint32, skillId uint32, skillLevel byte, attacked byte, damaged byte, attackedAndDamaged byte, stance byte, direction byte, rangedDirection byte, charge uint32, display byte, ranged bool, magic bool, speed byte, allDamage map[uint32][]uint32, x int16, y int16) {
      c := &attackCommand{
         WorldId:                  worldId,
         ChannelId:                channelId,
         MapId:                    mapId,
         CharacterId:              characterId,
         NumberAttacked:           attacked,
         NumberDamaged:            damaged,
         NumberAttackedAndDamaged: attackedAndDamaged,
         SkillId:                  skillId,
         SkillLevel:               skillLevel,
         Stance:                   stance,
         Direction:                direction,
         RangedDirection:          rangedDirection,
         Charge:                   charge,
         Display:                  display,
         Ranged:                   ranged,
         Magic:                    magic,
         Speed:                    speed,
         AllDamage:                allDamage,
         X:                        x,
         Y:                        y,
      }
      producer(CreateKey(int(characterId)), c)
   }
}
