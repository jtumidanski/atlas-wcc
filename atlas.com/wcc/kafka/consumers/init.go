package consumers

import (
	"atlas-wcc/kafka/handler"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	consumerGroupFormat = "World Channel Coordinator %d %d"

	EnableActionsCommand           = "enable_actions_command"
	ChangeMapEvent                 = "change_map_event"
	MapCharacterEvent              = "map_character_event"
	ControlMonsterEvent            = "control_monster_event"
	MonsterEvent                   = "monster_event"
	MonsterMovementEvent           = "monster_movement_event"
	CharacterMovementEvent         = "character_movement_event"
	CharacterMapMessageEvent       = "character_map_message_event"
	ExpressionChangedEvent         = "expression_changed_event"
	CharacterCreatedEvent          = "character_created_event"
	CharacterExperienceEvent       = "character_experience_event"
	InventoryModificationEvent     = "inventory_modification_event"
	CharacterLevelEvent            = "character_level_event"
	MesoGainedEvent                = "meso_gained_event"
	PickedUpItemEvent              = "picked_up_item_event"
	PickedUpNXEvent                = "picked_up_nx_event"
	DropReservationEvent           = "drop_reservation_event"
	PickupDropEvent                = "pickup_drop_event"
	DropEvent                      = "drop_event"
	CharacterSkillUpdateEvent      = "character_skill_update_event"
	CharacterStatisticEvent        = "character_statistic_event"
	ServerNoticeCommand            = "server_notice_command"
	MonsterKilledEvent             = "monster_killed_event"
	ShowDamageCharacterCommand     = "show_damage_character_command"
	NPCTalkCommand                 = "npc_talk_command"
	NPCTalkNumCommand              = "npc_talk_num_command"
	NPCTalkStyleCommand            = "npc_talk_style_command"
	DropExpireEvent                = "drop_expire_event"
	CharacterBuffEvent             = "character_buff_event"
	CharacterCancelBuffEvent       = "character_cancel_buff_event"
	CharacterEquipmentChangedEvent = "character_equipment_changed_event"
	InventoryFullCommand           = "inventory_full_command"
	CloseRangeAttackEvent          = "close_range_attack_event"
	RangeAttackEvent               = "range_attack_event"
	MagicAttackEvent               = "magic_attack_event"
	CharacterMPEaterEvent          = "character_mp_eater_event"
	ReactorStatusEvent             = "reactor_status_event"
	PartyStatusEvent               = "party_status_event"
	PartyMemberStatusEvent         = "party_member_status_event"
)

func CreateEventConsumers(l *logrus.Logger, wid byte, cid byte, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor ChannelEventProcessor) {
		createEventConsumer(l, wid, cid, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_ENABLE_ACTIONS", EnableActionsCommand, EnableActionsEventCreator(), HandleEnableActionsEvent())
	cec("TOPIC_CHANGE_MAP_EVENT", ChangeMapEvent, ChangeMapEventCreator(), HandleChangeMapEvent())
	cec("TOPIC_MAP_CHARACTER_EVENT", MapCharacterEvent, MapCharacterEventCreator(), HandleMapCharacterEvent())
	cec("TOPIC_CONTROL_MONSTER_EVENT", ControlMonsterEvent, MonsterControlEventCreator(), HandleMonsterControlEvent())
	cec("TOPIC_MONSTER_EVENT", MonsterEvent, MonsterEventCreator(), HandleMonsterEvent())
	cec("TOPIC_MONSTER_MOVEMENT", MonsterMovementEvent, MonsterMovementEventCreator(), HandleMonsterMovementEvent())
	cec("TOPIC_CHARACTER_MOVEMENT", CharacterMovementEvent, CharacterMovementEventCreator(), HandleCharacterMovementEvent())
	cec("TOPIC_CHARACTER_MAP_MESSAGE_EVENT", CharacterMapMessageEvent, CharacterMapMessageEventCreator(), HandleCharacterMapMessageEvent())
	cec("EXPRESSION_CHANGED", ExpressionChangedEvent, CharacterExpressionChangedEventCreator(), HandleCharacterExpressionChangedEvent())
	cec("TOPIC_CHARACTER_CREATED_EVENT", CharacterCreatedEvent, CharacterCreatedEventCreator(), HandleCharacterCreatedEvent())
	cec("TOPIC_CHARACTER_EXPERIENCE_EVENT", CharacterExperienceEvent, CharacterExperienceEventCreator(), HandleCharacterExperienceEvent())
	cec("TOPIC_INVENTORY_MODIFICATION", InventoryModificationEvent, CharacterInventoryModificationEventCreator(), HandleCharacterInventoryModificationEvent())
	cec("TOPIC_CHARACTER_LEVEL_EVENT", CharacterLevelEvent, CharacterLevelEventCreator(), HandleCharacterLevelEvent())
	cec("TOPIC_MESO_GAINED", MesoGainedEvent, CharacterMesoEventCreator(), HandleCharacterMesoEvent())
	cec("TOPIC_PICKED_UP_ITEM", PickedUpItemEvent, ItemPickedUpEventCreator(), HandleItemPickedUpEvent())
	cec("TOPIC_PICKED_UP_NX", PickedUpNXEvent, NXPickedUpEventCreator(), HandleNXPickedUpEvent())
	cec("TOPIC_DROP_RESERVATION_EVENT", DropReservationEvent, DropReservationEventCreator(), HandleDropReservationEvent())
	cec("TOPIC_PICKUP_DROP_EVENT", PickupDropEvent, DropPickedUpEventCreator(), HandleDropPickedUpEvent())
	cec("TOPIC_DROP_EVENT", DropEvent, DropEventCreator(), HandleDropEvent())
	cec("TOPIC_CHARACTER_SKILL_UPDATE_EVENT", CharacterSkillUpdateEvent, CharacterSkillUpdateEventCreator(), HandleCharacterSkillUpdateEvent())
	cec("TOPIC_CHARACTER_STAT_EVENT", CharacterStatisticEvent, CharacterStatUpdateEventCreator(), HandleCharacterStatUpdateEvent())
	cec("TOPIC_SERVER_NOTICE_COMMAND", ServerNoticeCommand, ServerNoticeEventCreator(), HandleServerNoticeEvent())
	cec("TOPIC_MONSTER_KILLED_EVENT", MonsterKilledEvent, MonsterKilledEventCreator(), HandleMonsterKilledEvent())
	cec("TOPIC_SHOW_DAMAGE_CHARACTER_COMMAND", ShowDamageCharacterCommand, CharacterDamagedEventCreator(), HandleCharacterDamagedEvent())
	cec("TOPIC_NPC_TALK_COMMAND", NPCTalkCommand, NPCTalkEventCreator(), HandleNPCTalkEvent())
	cec("TOPIC_NPC_TALK_NUM_COMMAND", NPCTalkNumCommand, EmptyNPCTalkNumCommandCreator(), HandleNPCTalkNumCommand())
	cec("TOPIC_NPC_TALK_STYLE_COMMAND", NPCTalkStyleCommand, EmptyNPCTalkStyleCommandCreator(), HandleNPCTalkStyleCommand())
	cec("TOPIC_DROP_EXPIRE_EVENT", DropExpireEvent, DropExpireEventCreator(), HandleDropExpireEvent())
	cec("TOPIC_CHARACTER_BUFF", CharacterBuffEvent, CharacterBuffEventCreator(), HandleCharacterBuffEvent())
	cec("TOPIC_CHARACTER_CANCEL_BUFF", CharacterCancelBuffEvent, CharacterCancelBuffEventCreator(), HandleCharacterCancelBuffEvent())
	cec("TOPIC_CHARACTER_EQUIP_CHANGED", CharacterEquipmentChangedEvent, CharacterEquipChangedEventCreator(), HandleCharacterEquipChangedEvent())
	cec("TOPIC_INVENTORY_FULL", InventoryFullCommand, InventoryFullCommandCreator(), HandleInventoryFullCommand())
	cec("TOPIC_CLOSE_RANGE_ATTACK_EVENT", CloseRangeAttackEvent, EmptyCloseRangeAttackEventCreator(), HandleCloseRangeAttackEvent())
	cec("TOPIC_RANGE_ATTACK_EVENT", RangeAttackEvent, EmptyRangeAttackEventCreator(), HandleRangeAttackEvent())
	cec("TOPIC_MAGIC_ATTACK_EVENT", MagicAttackEvent, EmptyMagicAttackEventCreator(), HandleMagicAttackEvent())
	cec("TOPIC_CHARACTER_MP_EATER_EVENT", CharacterMPEaterEvent, EmptyMPEaterEventCreator(), HandleMPEaterEvent())
	cec("TOPIC_REACTOR_STATUS_EVENT", ReactorStatusEvent, EmptyReactorStatusEventCreator(), HandleReactorStatusEvent())
	cec("TOPIC_PARTY_STATUS", PartyStatusEvent, EmptyPartyStatusEventCreator(), HandlePartyStatusEvent())
	cec("TOPIC_PARTY_MEMBER_STATUS", PartyMemberStatusEvent, EmptyPartyMemberStatusEventCreator(), HandlePartyMemberStatusEvent())
}

func createEventConsumer(l *logrus.Logger, wid byte, cid byte, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor ChannelEventProcessor) {
	wg.Add(1)
	h := func(logger logrus.FieldLogger, span opentracing.Span, event interface{}) {
		processor(logger, span, wid, cid, event)
	}
	groupId := fmt.Sprintf(consumerGroupFormat, wid, cid)
	go NewConsumer(l, ctx, wg, name, topicToken, groupId, emptyEventCreator, h)
}
