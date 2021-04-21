package consumers

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	consumerGroupFormat = "World Channel Coordinator %d %d"
)

func CreateEventConsumers(l *logrus.Logger, wid byte, cid byte) {
	cec := func(topicToken string, emptyEventCreator EmptyEventCreator, processor ChannelEventProcessor) {
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
	cec("TOPIC_DROP_EXPIRE_EVENT", DropExpireEventCreator(), HandleDropExpireEvent())
}

func createEventConsumer(l *logrus.Logger, wid byte, cid byte, topicToken string, emptyEventCreator EmptyEventCreator, processor ChannelEventProcessor) {
	groupId := fmt.Sprintf(consumerGroupFormat, wid, cid)

	h := func(logger logrus.FieldLogger, event interface{}) {
		processor(logger, wid, cid, event)
	}

	c := NewConsumer(l, context.Background(), h,
		SetGroupId(groupId),
		SetTopicToken(topicToken),
		SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}