package _map

import (
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/drop"
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/monster"
	"atlas-wcc/npc"
	"atlas-wcc/reactor"
	"atlas-wcc/session"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	consumerNameCharacterDamage           = "show_damage_character_command"
	consumerNameMessage                   = "character_map_message_event"
	consumerNameMapCharacter              = "map_character_event"
	consumerNameCharacterLevel            = "character_level_event"
	consumerNameMapChange                 = "change_map_event"
	consumerNameMPEater                   = "character_mp_eater_event"
	consumerNameCharacterExpressionChange = "expression_changed_event"
	consumerNameCharacterMovement         = "character_movement_event"
	consumerNameCharacterEquipmentChanged = "character_equipment_changed_event"
	consumerNameDropEvent                 = "drop_event"
	consumerNameDropExpire                = "drop_expire_event"
	consumerNamePickup                    = "pickup_drop_event"
	consumerNameMonsterEvent              = "monster_event"
	consumerNameMonsterMovement           = "monster_movement_event"
	consumerNameMonsterKilled             = "monster_killed_event"
	consumerNameReactorStatus             = "reactor_status_event"
	topicTokenMessage                     = "TOPIC_CHARACTER_MAP_MESSAGE_EVENT"
	topicTokenMapCharacter                = "TOPIC_MAP_CHARACTER_EVENT"
	topicTokenCharacterLevel              = "TOPIC_CHARACTER_LEVEL_EVENT"
	topicTokenMapChange                   = "TOPIC_CHANGE_MAP_EVENT"
	topicTokenMPEater                     = "TOPIC_CHARACTER_MP_EATER_EVENT"
	topicTokenCharacterDamage             = "TOPIC_SHOW_DAMAGE_CHARACTER_COMMAND"
	topicTokenCharacterExpressionChange   = "EXPRESSION_CHANGED"
	topicTokenCharacterMovement           = "TOPIC_CHARACTER_MOVEMENT"
	topicTokenCharacterEquipmentChanged   = "TOPIC_CHARACTER_EQUIP_CHANGED"
	topicTokenDropEvent                   = "TOPIC_DROP_EVENT"
	topicTokenDropExpire                  = "TOPIC_DROP_EXPIRE_EVENT"
	topicTokenPickup                      = "TOPIC_PICKUP_DROP_EVENT"
	topicTokenMonsterEvent                = "TOPIC_MONSTER_EVENT"
	topicTokenMonsterMovement             = "TOPIC_MONSTER_MOVEMENT"
	topicTokenMonsterKilled               = "TOPIC_MONSTER_KILLED_EVENT"
	topicTokenReactorStatus               = "TOPIC_REACTOR_STATUS_EVENT"
)

func MessageConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[messageEvent](consumerNameMessage, topicTokenMessage, groupId, handleMessage(wid, cid))
	}
}

type messageEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Message     string `json:"message"`
	GM          bool   `json:"gm"`
	Show        bool   `json:"show"`
}

func handleMessage(wid byte, cid byte) kafka.HandlerFunc[messageEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event messageEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, session.Announce(WriteChatText(l)(event.CharacterId, event.Message, event.GM, event.Show)))
	}
}

func MapCharacterConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[mapCharacterEvent](consumerNameMapCharacter, topicTokenMapCharacter, groupId, HandleMapCharacterEvent(wid, cid))
	}
}

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func HandleMapCharacterEvent(wid byte, cid byte) kafka.HandlerFunc[mapCharacterEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event mapCharacterEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		if event.Type == "ENTER" {
			session.ForSessionByCharacterId(event.CharacterId, enterMap(l, span)(event))
		} else if event.Type == "EXIT" {
			ForOtherSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, event.CharacterId, session.Announce(WriteRemoveCharacterFromMap(l)(event.CharacterId)))
		} else {
			l.Warnf("Received a unhandled map character event type of %s.", event.Type)
			return
		}
	}
}

func enterMap(l logrus.FieldLogger, span opentracing.Span) func(event mapCharacterEvent) model.Operator[session.Model] {
	return func(event mapCharacterEvent) model.Operator[session.Model] {
		return func(s session.Model) error {
			ids, err := GetCharacterIdsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId)
			if err != nil {
				l.WithError(err).Errorf("No characters found in map %d for world %d and channel %d.", event.MapId, event.WorldId, event.ChannelId)
				return err
			}

			cm := make(map[uint32]character.Model)
			for _, cId := range ids {
				c, err := character.GetCharacterById(l, span)(cId)
				if err != nil {
					//log something
				} else {
					cm[c.Attributes().Id()] = c
				}
			}

			// Spawn new character for other character.
			for k, v := range cm {
				if k != event.CharacterId {
					session.ForSessionByCharacterId(k, session.Announce(WriteSpawnCharacter(l)(v, cm[event.CharacterId], true)))
				}
			}

			// Spawn other characters for incoming character.
			for k, v := range cm {
				if k != event.CharacterId {
					err := session.Announce(WriteSpawnCharacter(l)(cm[event.CharacterId], v, false))(s)
					if err != nil {
						l.WithError(err).Errorf("Unable to spawn character %d for %d", v.Attributes().Id(), event.CharacterId)
					}
				}
			}

			// Spawn NPCs for incoming character.

			npc.ForEachInMap(l, span)(event.MapId, npc.SpawnNPCForSession(l)(s))

			// Spawn monsters for incoming character.
			monster.ForEachInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, monster.SpawnForSession(l)(s))

			// Spawn drops for incoming character.
			drop.ForEachInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, drop.SpawnDropForSession(l)(s))

			// Spawn reactors for incoming character.
			reactor.ForEachAliveInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, reactor.SpawnForSession(l)(s))
			return err
		}
	}
}

func CharacterLevelConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterLevelEvent](consumerNameCharacterLevel, topicTokenCharacterLevel, groupId, handleCharacterLevel(wid, cid))
	}
}

type characterLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func handleCharacterLevel(wid byte, cid byte) kafka.HandlerFunc[characterLevelEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterLevelEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		c, err := properties.GetById(l, span)(event.CharacterId)
		if err != nil {
			return
		}

		ForOtherSessionsInMap(l, span)(wid, cid, c.MapId(), c.Id(), session.Announce(WriteShowForeignEffect(l)(event.CharacterId, 0)))
	}
}

func MapChangeConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[mapChangedEvent](consumerNameMapChange, topicTokenMapChange, groupId, handleChangeMap(wid, cid))
	}
}

type mapChangedEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func handleChangeMap(wid byte, cid byte) kafka.HandlerFunc[mapChangedEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event mapChangedEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, warpCharacter(l, span, event))
	}
}

func warpCharacter(l logrus.FieldLogger, span opentracing.Span, event mapChangedEvent) model.Operator[session.Model] {
	catt, err := properties.GetById(l, span)(event.CharacterId)
	if err != nil {
		l.WithError(err).Errorf("Unable to retrieve character %d properties", event.CharacterId)
		return model.ErrorOperator[session.Model](err)
	}
	return session.Announce(WriteWarpToMap(l)(event.ChannelId, event.MapId, event.PortalId, catt.Hp()))
}

func MPEaterConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[mpEaterEvent](consumerNameMPEater, topicTokenMPEater, groupId, handleMPEater(wid, cid))
	}
}

type mpEaterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func handleMPEater(_ byte, _ byte) kafka.HandlerFunc[mpEaterEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event mpEaterEvent) {
		session.ForSessionByCharacterId(event.CharacterId, session.Announce(WriteShowOwnBuff(l)(1, event.SkillId)))
		ForOtherSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, event.CharacterId, session.Announce(WriteShowBuffEffect(l)(event.CharacterId, 1, event.SkillId, 3)))
	}
}

func CharacterDamageConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterDamageEvent](consumerNameCharacterDamage, topicTokenCharacterDamage, groupId, handleCharacterDamage(wid, cid))
	}
}

type characterDamageEvent struct {
	CharacterId     uint32 `json:"characterId"`
	MapId           uint32 `json:"mapId"`
	MonsterId       uint32 `json:"monsterId"`
	MonsterUniqueId uint32 `json:"monsterUniqueId"`
	SkillId         int8   `json:"skillId"`
	Damage          int32  `json:"damage"`
	Fake            uint32 `json:"fake"`
	Direction       int8   `json:"direction"`
	X               int16  `json:"x"`
	Y               int16  `json:"y"`
	PGMR            bool   `json:"pgmr"`
	PGMR1           byte   `json:"pgmr1"`
	PG              bool   `json:"pg"`
}

func handleCharacterDamage(wid byte, cid byte) kafka.HandlerFunc[characterDamageEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterDamageEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, session.Announce(WriteCharacterDamaged(l)(event.SkillId, event.MonsterId, event.CharacterId, event.Damage, event.Fake, event.Direction, event.PGMR, event.PGMR1, event.PG, event.MonsterUniqueId, event.X, event.Y)))
	}
}

func CharacterExpressionChangeConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterExpressionChangedEvent](consumerNameCharacterExpressionChange, topicTokenCharacterExpressionChange, groupId, handleCharacterExpressionChange(wid, cid))
	}
}

type characterExpressionChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Expression  uint32 `json:"expression"`
}

func handleCharacterExpressionChange(wid byte, cid byte) kafka.HandlerFunc[characterExpressionChangedEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterExpressionChangedEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, session.Announce(WriteCharacterExpression(l)(event.CharacterId, event.Expression)))
	}
}

func CharacterMovementConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterMovementEvent](consumerNameCharacterMovement, topicTokenCharacterMovement, groupId, handleCharacterMovement(wid, cid))
	}
}

type characterMovementEvent struct {
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

func handleCharacterMovement(wid byte, cid byte) kafka.HandlerFunc[characterMovementEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterMovementEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		c, err := properties.GetById(l, span)(event.CharacterId)
		if err != nil {
			return
		}

		ForOtherSessionsInMap(l, span)(wid, cid, c.MapId(), c.Id(), session.Announce(WriteMoveCharacter(l)(event.CharacterId, event.RawMovement)))
	}
}

func CharacterEquipmentChangedConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterEquipmentChangedEvent](consumerNameCharacterEquipmentChanged, topicTokenCharacterEquipmentChanged, groupId, handleCharacterEquipmentChanged(wid, cid))
	}
}

type characterEquipmentChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Change      string `json:"change"`
}

func handleCharacterEquipmentChanged(wid byte, cid byte) kafka.HandlerFunc[characterEquipmentChangedEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterEquipmentChangedEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		c, err := properties.GetById(l, span)(event.CharacterId)
		if err != nil {
			return
		}

		ForOtherSessionsInMap(l, span)(wid, cid, c.MapId(), c.Id(), updateCharacterAppearance(l, span)(event.CharacterId))
	}
}

func updateCharacterAppearance(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) model.Operator[session.Model] {
	return func(characterId uint32) model.Operator[session.Model] {
		return func(s session.Model) error {
			r, err := character.GetCharacterById(l, span)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve character %d details.", s.CharacterId())
				return err
			}
			c, err := character.GetCharacterById(l, span)(characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve character %d details.", characterId)
				return err
			}
			err = session.Announce(WriteCharacterLookUpdated(l)(r, c))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to %d that character %d has changed their look.", s.CharacterId(), characterId)
			}
			return err
		}
	}
}

func DropEventConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[dropEvent](consumerNameDropEvent, topicTokenDropEvent, groupId, handleDropEvent(wid, cid))
	}
}

type dropEvent struct {
	WorldId         byte   `json:"worldId"`
	ChannelId       byte   `json:"channelId"`
	MapId           uint32 `json:"mapId"`
	UniqueId        uint32 `json:"uniqueId"`
	ItemId          uint32 `json:"itemId"`
	Quantity        uint32 `json:"quantity"`
	Meso            uint32 `json:"meso"`
	DropType        byte   `json:"dropType"`
	DropX           int16  `json:"dropX"`
	DropY           int16  `json:"dropY"`
	OwnerId         uint32 `json:"ownerId"`
	OwnerPartyId    uint32 `json:"ownerPartyId"`
	DropTime        uint64 `json:"dropTime"`
	DropperUniqueId uint32 `json:"dropperUniqueId"`
	DropperX        int16  `json:"dropperX"`
	DropperY        int16  `json:"dropperY"`
	PlayerDrop      bool   `json:"playerDrop"`
	Mod             byte   `json:"mod"`
}

func handleDropEvent(wid byte, cid byte) kafka.HandlerFunc[dropEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event dropEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, showDrop(l, event))
	}
}

func showDrop(l logrus.FieldLogger, event dropEvent) model.Operator[session.Model] {
	return func(s session.Model) error {
		a := uint32(0)
		if event.ItemId != 0 {
			a = 0
		} else {
			a = event.Meso
		}
		err := session.Announce(drop.WriteDropItemFromMapObject(l)(event.UniqueId, event.ItemId, event.Meso, a,
			event.DropperUniqueId, event.DropType, event.OwnerId, event.OwnerPartyId, s.CharacterId(), 0,
			event.DropTime, event.DropX, event.DropY, event.DropperX, event.DropperY, event.PlayerDrop, event.Mod))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write drop in map for character %d", s.CharacterId())
		}
		return err
	}
}

func ExpireDropConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[dropExpireEvent](consumerNameDropExpire, topicTokenDropExpire, groupId, handleDropExpiration(wid, cid))
	}
}

type dropExpireEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
}

func handleDropExpiration(wid byte, cid byte) kafka.HandlerFunc[dropExpireEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event dropExpireEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, session.Announce(drop.WriteRemoveItem(l)(event.UniqueId, 0, 0)))
	}
}

func PickupConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[pickupEvent](consumerNamePickup, topicTokenPickup, groupId, handlePickup(wid, cid))
	}
}

type pickupEvent struct {
	DropId      uint32 `json:"dropId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
}

func handlePickup(wid byte, cid byte) kafka.HandlerFunc[pickupEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event pickupEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, session.Announce(drop.WriteRemoveItem(l)(event.DropId, 2, event.CharacterId)))
	}
}

func MonsterEventConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[monsterEvent](consumerNameMonsterEvent, topicTokenMonsterEvent, groupId, handleMonsterEvent(wid, cid))
	}
}

type monsterEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
	MonsterId uint32 `json:"monsterId"`
	ActorId   uint32 `json:"actorId"`
	Type      string `json:"type"`
}

func handleMonsterEvent(wid byte, cid byte) kafka.HandlerFunc[monsterEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event monsterEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		m, err := monster.GetById(l, span)(event.UniqueId)
		if err != nil {
			l.WithError(err).Errorf("Unable to monster %d to create.", event.UniqueId)
			return
		}

		var h model.Operator[session.Model]
		if event.Type == "CREATED" {
			h = session.Announce(monster.WriteSpawnMonster(l)(m, false))
		} else if event.Type == "DESTROYED" {
			h = session.Announce(monster.WriteKillMonster(l)(m.UniqueId(), false), monster.WriteKillMonster(l)(m.UniqueId(), true))
		} else {
			l.Warnf("Unable to handle %s event type for monster events.", event.Type)
			return
		}

		ForSessionsInMap(l, span)(wid, cid, event.MapId, h)
	}
}

func MonsterMovementConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[monsterMovementEvent](consumerNameMonsterMovement, topicTokenMonsterMovement, groupId, handleMovement(wid, cid))
	}
}

type monsterMovementEvent struct {
	UniqueId      uint32      `json:"uniqueId"`
	ObserverId    uint32      `json:"observerId"`
	SkillPossible bool        `json:"skillPossible"`
	Skill         int8        `json:"skill"`
	SkillId       uint32      `json:"skillId"`
	SkillLevel    uint32      `json:"skillLevel"`
	Option        uint16      `json:"option"`
	StartX        int16       `json:"startX"`
	StartY        int16       `json:"startY"`
	EndX          int16       `json:"endX"`
	EndY          int16       `json:"endY"`
	Stance        byte        `json:"stance"`
	RawMovement   rawMovement `json:"rawMovement"`
}

func handleMovement(wid byte, cid byte) kafka.HandlerFunc[monsterMovementEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event monsterMovementEvent) {
		if _, err := session.GetByCharacterId(event.ObserverId); err != nil {
			return
		}

		c, err := properties.GetById(l, span)(event.ObserverId)
		if err != nil {
			return
		}

		ForOtherSessionsInMap(l, span)(wid, cid, c.MapId(), c.Id(), session.Announce(monster.WriteMoveMonster(l)(event.UniqueId,
			event.SkillPossible, event.Skill, event.SkillId, event.SkillLevel, event.Option, event.StartX,
			event.StartY, event.RawMovement)))
	}
}

func MonsterDeathConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[killedEvent](consumerNameMonsterKilled, topicTokenMonsterKilled, groupId, handleDeath(wid, cid))
	}
}

type entry struct {
	CharacterId uint32 `json:"characterId"`
	Damage      uint64 `json:"damage"`
}

type killedEvent struct {
	WorldId       byte    `json:"worldId"`
	ChannelId     byte    `json:"channelId"`
	MapId         uint32  `json:"mapId"`
	UniqueId      uint32  `json:"uniqueId"`
	MonsterId     uint32  `json:"monsterId"`
	X             int16   `json:"x"`
	Y             int16   `json:"y"`
	KillerId      uint32  `json:"killerId"`
	DamageEntries []entry `json:"damageEntries"`
}

func handleDeath(wid byte, cid byte) kafka.HandlerFunc[killedEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event killedEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		l.Infof("Character %d killed %d.", event.UniqueId, event.KillerId)
		ForSessionsInMap(l, span)(wid, cid, event.MapId, session.Announce(monster.WriteKillMonster(l)(event.UniqueId, true)))
	}
}

func ReactorStatusConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[statusEvent](consumerNameReactorStatus, topicTokenReactorStatus, groupId, handleStatus(wid, cid))
	}
}

type statusEvent struct {
	WorldId   byte   `json:"world_id"`
	ChannelId byte   `json:"channel_id"`
	MapId     uint32 `json:"map_id"`
	Id        uint32 `json:"id"`
	Status    string `json:"status"`
	Stance    uint16 `json:"stance"`
}

func handleStatus(wid byte, cid byte) kafka.HandlerFunc[statusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statusEvent) {
		if event.WorldId != wid || event.ChannelId != cid {
			return
		}
		if event.Status == "CREATED" {
			ForSessionsInMap(l, span)(wid, cid, event.MapId, reactor.CreateForSession(l, span)(event.Id, event.Stance))
		} else if event.Status == "TRIGGERED" {
			ForSessionsInMap(l, span)(wid, cid, event.MapId, reactor.HitForSession(l, span)(event.Id, event.Stance))
		} else if event.Status == "DESTROYED" {
			ForSessionsInMap(l, span)(wid, cid, event.MapId, reactor.DestroyForSession(l, span)(event.Id))
		}
	}
}
