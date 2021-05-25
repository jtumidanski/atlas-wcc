package consumers

import (
	"atlas-wcc/kafka/handler"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	consumerGroupFormat = "World Channel Coordinator %d %d"
)

func CreateEventConsumers(l *logrus.Logger, wid byte, cid byte) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor ChannelEventProcessor) {
		createEventConsumer(l, wid, cid, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_ENABLE_ACTIONS", EnableActionsEventCreator(), HandleEnableActionsEvent())
	cec("TOPIC_CHANGE_MAP_EVENT", ChangeMapEventCreator(), HandleChangeMapEvent())
	cec("TOPIC_MAP_CHARACTER_EVENT", MapCharacterEventCreator(), HandleMapCharacterEvent())
	cec("TOPIC_CONTROL_MONSTER_EVENT", MonsterControlEventCreator(), HandleMonsterControlEvent())
	cec("TOPIC_MONSTER_EVENT", MonsterEventCreator(), HandleMonsterEvent())
	cec("TOPIC_MONSTER_MOVEMENT", MonsterMovementEventCreator(), HandleMonsterMovementEvent())
	cec("TOPIC_CHARACTER_MOVEMENT", CharacterMovementEventCreator(), HandleCharacterMovementEvent())
	cec("TOPIC_CHARACTER_MAP_MESSAGE_EVENT", CharacterMapMessageEventCreator(), HandleCharacterMapMessageEvent())
	cec("EXPRESSION_CHANGED", CharacterExpressionChangedEventCreator(), HandleCharacterExpressionChangedEvent())
	cec("TOPIC_CHARACTER_CREATED_EVENT", CharacterCreatedEventCreator(), HandleCharacterCreatedEvent())
	cec("TOPIC_CHARACTER_EXPERIENCE_EVENT", CharacterExperienceEventCreator(), HandleCharacterExperienceEvent())
	cec("TOPIC_INVENTORY_MODIFICATION", CharacterInventoryModificationEventCreator(), HandleCharacterInventoryModificationEvent())
	cec("TOPIC_CHARACTER_LEVEL_EVENT", CharacterLevelEventCreator(), HandleCharacterLevelEvent())
	cec("TOPIC_MESO_GAINED", CharacterMesoEventCreator(), HandleCharacterMesoEvent())
	cec("TOPIC_PICKED_UP_ITEM", ItemPickedUpEventCreator(), HandleItemPickedUpEvent())
	cec("TOPIC_PICKED_UP_NX", NXPickedUpEventCreator(), HandleNXPickedUpEvent())
	cec("TOPIC_DROP_RESERVATION_EVENT", DropReservationEventCreator(), HandleDropReservationEvent())
	cec("TOPIC_PICKUP_DROP_EVENT", DropPickedUpEventCreator(), HandleDropPickedUpEvent())
	cec("TOPIC_DROP_EVENT", DropEventCreator(), HandleDropEvent())
	cec("TOPIC_CHARACTER_SKILL_UPDATE_EVENT", CharacterSkillUpdateEventCreator(), HandleCharacterSkillUpdateEvent())
	cec("TOPIC_CHARACTER_STAT_EVENT", CharacterStatUpdateEventCreator(), HandleCharacterStatUpdateEvent())
	cec("TOPIC_SERVER_NOTICE_COMMAND", ServerNoticeEventCreator(), HandleServerNoticeEvent())
	cec("TOPIC_MONSTER_KILLED_EVENT", MonsterKilledEventCreator(), HandleMonsterKilledEvent())
	cec("TOPIC_SHOW_DAMAGE_CHARACTER_COMMAND", CharacterDamagedEventCreator(), HandleCharacterDamagedEvent())
	cec("TOPIC_NPC_TALK_COMMAND", NPCTalkEventCreator(), HandleNPCTalkEvent())
	cec("TOPIC_NPC_TALK_NUM_COMMAND", EmptyNPCTalkNumCommandCreator(), HandleNPCTalkNumCommand())
	cec("TOPIC_NPC_TALK_STYLE_COMMAND", EmptyNPCTalkStyleCommandCreator(), HandleNPCTalkStyleCommand())
	cec("TOPIC_DROP_EXPIRE_EVENT", DropExpireEventCreator(), HandleDropExpireEvent())
	cec("TOPIC_CHARACTER_BUFF", CharacterBuffEventCreator(), HandleCharacterBuffEvent())
	cec("TOPIC_CHARACTER_CANCEL_BUFF", CharacterCancelBuffEventCreator(), HandleCharacterCancelBuffEvent())
	cec("TOPIC_CHARACTER_EQUIP_CHANGED", CharacterEquipChangedEventCreator(), HandleCharacterEquipChangedEvent())
	cec("TOPIC_INVENTORY_FULL", InventoryFullCommandCreator(), HandleInventoryFullCommand())
	cec("TOPIC_CLOSE_RANGE_ATTACK_EVENT", EmptyCloseRangeAttackEventCreator(), HandleCloseRangeAttackEvent())
	cec("TOPIC_RANGE_ATTACK_EVENT", EmptyRangeAttackEventCreator(), HandleRangeAttackEvent())
	cec("TOPIC_MAGIC_ATTACK_EVENT", EmptyMagicAttackEventCreator(), HandleMagicAttackEvent())
	cec("TOPIC_CHARACTER_MP_EATER_EVENT", EmptyMPEaterEventCreator(), HandleMPEaterEvent())
}

func createEventConsumer(l *logrus.Logger, wid byte, cid byte, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor ChannelEventProcessor) {
	groupId := fmt.Sprintf(consumerGroupFormat, wid, cid)
	h := func(logger logrus.FieldLogger, event interface{}) {
		processor(logger, wid, cid, event)
	}
	go NewConsumer(l, topicToken, groupId, emptyEventCreator, h)
}
