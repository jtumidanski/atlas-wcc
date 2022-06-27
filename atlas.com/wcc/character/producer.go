package character

import (
	"atlas-wcc/kafka"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strings"
)

type healthAdjustmentEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func AdjustHealth(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, amount int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ADJUST_HEALTH")
	return func(characterId uint32, amount int16) {
		e := &healthAdjustmentEvent{
			CharacterId: characterId,
			Amount:      amount,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type manaAdjustmentEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func AdjustMana(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, amount int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ADJUST_MANA")
	return func(characterId uint32, amount int16) {
		e := &manaAdjustmentEvent{
			CharacterId: characterId,
			Amount:      amount,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

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

func Attack(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, skillId uint32, skillLevel byte, attacked byte, damaged byte, attackedAndDamaged byte, stance byte, direction byte, rangedDirection byte, charge uint32, display byte, ranged bool, magic bool, speed byte, allDamage map[uint32][]uint32, x int16, y int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_ATTACK_COMMAND")
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
		producer(kafka.CreateKey(int(characterId)), c)
	}
}

type damageCommand struct {
	CharacterId     uint32 `json:"characterId"`
	MonsterId       uint32 `json:"monsterId"`
	MonsterUniqueId uint32 `json:"monsterUniqueId"`
	DamageFrom      int8   `json:"damageFrom"`
	Element         byte   `json:"element"`
	Damage          int32  `json:"damage"`
	Direction       int8   `json:"direction"`
}

func Damage(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, monsterIdFrom uint32, uniqueId uint32, damageFrom int8, element byte, damage int32, direction int8) {
	producer := kafka.ProduceEvent(l, span, "DAMAGE_CHARACTER")
	return func(characterId uint32, monsterIdFrom uint32, uniqueId uint32, damageFrom int8, element byte, damage int32, direction int8) {
		e := &damageCommand{
			CharacterId:     characterId,
			MonsterId:       monsterIdFrom,
			MonsterUniqueId: uniqueId,
			DamageFrom:      damageFrom,
			Element:         element,
			Damage:          damage,
			Direction:       direction,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type distributeApCommand struct {
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func DistributeAp(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, attributeType string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ASSIGN_AP_COMMAND")
	return func(characterId uint32, attributeType string) {
		e := &distributeApCommand{
			CharacterId: characterId,
			Type:        attributeType,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type distributeSpCommand struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func DistributeSp(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, skillId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ASSIGN_SP_COMMAND")
	return func(characterId uint32, skillId uint32) {
		e := &distributeSpCommand{
			CharacterId: characterId,
			SkillId:     skillId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type dropItemCommand struct {
	WorldId       byte   `json:"worldId"`
	ChannelId     byte   `json:"channelId"`
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Quantity      int16  `json:"quantity"`
}

func DropItem(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, inventoryType int8, source int16, quantity int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_DROP_ITEM")
	return func(worldId byte, channelId byte, characterId uint32, inventoryType int8, source int16, quantity int16) {
		e := &dropItemCommand{WorldId: worldId, ChannelId: channelId, CharacterId: characterId, InventoryType: inventoryType, Source: source, Quantity: quantity}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type equipItemCommand struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func EquipItem(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, source int16, destination int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_EQUIP_ITEM")
	return func(characterId uint32, source int16, destination int16) {
		e := &equipItemCommand{CharacterId: characterId, Source: source, Destination: destination}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type expressionChangeCommand struct {
	CharacterId uint32 `json:"characterId"`
	Emote       uint32 `json:"emote"`
}

func ChangeExpression(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, emote uint32) {
	producer := kafka.ProduceEvent(l, span, "CHANGE_FACIAL_EXPRESSION")
	return func(characterId uint32, emote uint32) {
		e := &expressionChangeCommand{
			CharacterId: characterId,
			Emote:       emote,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type mapMessageCommand struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Message     string `json:"message"`
	GM          bool   `json:"gm"`
	Show        bool   `json:"show"`
}

func SendMapMessage(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, message string, gm bool, show bool) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_MAP_MESSAGE_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, message string, gm bool, show bool) {
		e := &mapMessageCommand{
			CharacterId: characterId,
			MapId:       mapId,
			Message:     message,
			GM:          gm,
			Show:        show,
		}
		producer(kafka.CreateKey(int(worldId)*1000+int(channelId)), e)
	}
}

type moveItemCommand struct {
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Destination   int16  `json:"destination"`
}

func MoveItem(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MOVE_ITEM")
	return func(characterId uint32, inventoryType int8, source int16, destination int16) {
		e := &moveItemCommand{CharacterId: characterId, InventoryType: inventoryType, Source: source, Destination: destination}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type movementEvent struct {
	WorldId     byte        `json:"worldId"`
	ChannelId   byte        `json:"channelId"`
	CharacterId uint32      `json:"characterId"`
	X           int16       `json:"x"`
	Y           int16       `json:"y"`
	Stance      byte        `json:"stance"`
	RawMovement rawMovement `json:"rawMovement"`
}

type rawMovement []byte

func (m rawMovement) MarshalJSON() ([]byte, error) {
	var result string
	if m == nil {
		result = "[]"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", m)), ",")
	}
	return []byte(result), nil
}

func Move(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, x int16, y int16, stance byte, rawMovement []byte) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_MOVEMENT")
	return func(worldId byte, channelId byte, characterId uint32, x int16, y int16, stance byte, rawMovement []byte) {
		e := &movementEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
			X:           x,
			Y:           y,
			Stance:      stance,
			RawMovement: rawMovement,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type reserveDropCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func ReserveDrop(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, dropId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_RESERVE_DROP_COMMAND")
	return func(characterId uint32, dropId uint32) {
		e := &reserveDropCommand{
			CharacterId: characterId,
			DropId:      dropId,
		}
		producer(kafka.CreateKey(int(dropId)), e)
	}
}

type applySkillCommand struct {
	CharacterId uint32
	SkillId     uint32
	Level       uint8
	X           int16
	Y           int16
}

func ApplySkill(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, skillId uint32, level uint8, x int16, y int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_APPLY_SKILL_COMMAND")
	return func(characterId uint32, skillId uint32, level uint8, x int16, y int16) {
		e := &applySkillCommand{
			CharacterId: characterId,
			SkillId:     skillId,
			Level:       level,
			X:           x,
			Y:           y,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type MonsterMagnetData struct {
	MonsterId uint32
	Success   uint8
}

type applyMonsterMagnetCommand struct {
	CharacterId uint32
	SkillId     uint32
	Level       uint8
	Direction   int8
	Data        []MonsterMagnetData
}

func ApplyMonsterMagnet(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, skillId uint32, level uint8, direction int8, data []MonsterMagnetData) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_APPLY_MONSTER_MAGNET_COMMAND")
	return func(characterId uint32, skillId uint32, level uint8, direction int8, data []MonsterMagnetData) {
		e := applyMonsterMagnetCommand{
			CharacterId: characterId,
			SkillId:     skillId,
			Level:       level,
			Direction:   direction,
			Data:        data,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type unequipItemCommand struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func UnequipItem(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, source int16, destination int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_UNEQUIP_ITEM")
	return func(characterId uint32, source int16, destination int16) {
		e := &unequipItemCommand{CharacterId: characterId, Source: source, Destination: destination}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type adjustMesoEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int32  `json:"amount"`
	Show        bool   `json:"show"`
}

func emitMesoAdjustment(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, amount int32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ADJUST_MESO")
	return func(characterId uint32, amount int32) {
		event := &adjustMesoEvent{characterId, amount, true}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}
