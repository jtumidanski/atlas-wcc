package socket

import (
	"atlas-wcc/session"
	"atlas-wcc/socket/request"
	"atlas-wcc/socket/request/handler"
	"context"
	"github.com/jtumidanski/atlas-socket"
	request2 "github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
	"sync"
)

func CreateSocketService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) func(worldId byte, channelId byte, port int) {
	return func(worldId byte, channelId byte, port int) {
		go func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				wg.Add(1)
				defer wg.Done()
				err := socket.Run(l, handlerProducer(l),
					socket.SetPort(port),
					socket.SetSessionCreator(session.Create(l, session.Registry())(worldId, channelId)),
					socket.SetSessionMessageDecryptor(session.Decrypt(l, session.Registry())),
					socket.SetSessionDestroyer(session.DestroyByIdWithSpan(l, session.Registry())),
				)
				if err != nil {
					l.WithError(err).Errorf("Socket service encountered error")
				}
			}()

			<-ctx.Done()
			l.Infof("Shutting down server on port 8484")
		}()
	}
}

const (
	Pong                      = "pong"
	CharacterLoggedIn         = "character_logged_in"
	ChangeMapSpecial          = "change_map_special"
	MoveCharacter             = "move_character"
	ChangeMap                 = "change_map"
	MoveLife                  = "move_life"
	GeneralChat               = "general_chat"
	ChangeChannel             = "change_channel"
	CharacterExpression       = "character_expression"
	CharacterCloseRangeAttack = "character_close_range_attack"
	CharacterRangedAttack     = "character_ranged_attack"
	CharacterMagicAttack      = "character_magic_attack"
	CharacterDistributeAP     = "character_distribute_ap"
	CharacterDistributeSP     = "character_distribute_sp"
	CharacterHealOverTime     = "character_heal_over_time"
	CharacterItemPickUp       = "character_item_pick_up"
	NPCAction                 = "npc_action"
	NPCTalkMore               = "npc_talk_more"
	NPCTalk                   = "npc_talk"
	CharacterDamage           = "character_damage"
	MoveItem                  = "move_item"
	SpecialMove               = "special_move"
	QuestAction               = "quest_action"
	InnerPortal               = "inner_portal"
	ChangeKeyMap              = "change_key_map"
	ReactorHit                = "reactor_hit"
)

func handlerProducer(l logrus.FieldLogger) socket.MessageHandlerProducer {
	handlers := make(map[uint16]request2.Handler)
	hr := func(op uint16, name string, v request.MessageValidator, h request.MessageHandler) {
		handlers[op] = request.AdaptHandler(l, name, v, h)
	}

	hr(handler.OpCodePong, Pong, request.NoOpValidator, request.NoOpHandler)
	hr(handler.OpCharacterLoggedIn, CharacterLoggedIn, request.NoOpValidator, handler.CharacterLoggedInHandler)
	hr(handler.OpChangeMapSpecial, ChangeMapSpecial, request.LoggedInValidator, handler.ChangeMapSpecialHandler)
	hr(handler.OpMoveCharacter, MoveCharacter, request.LoggedInValidator, handler.MoveCharacterHandler)
	hr(handler.OpChangeMap, ChangeMap, request.LoggedInValidator, handler.ChangeMapHandler)
	hr(handler.OpMoveLife, MoveLife, request.LoggedInValidator, handler.MoveLifeHandler)
	hr(handler.OpGeneralChat, GeneralChat, request.LoggedInValidator, handler.GeneralChatHandler)
	hr(handler.OpChangeChannel, ChangeChannel, request.LoggedInValidator, handler.ChangeChannelHandler)
	hr(handler.OpCharacterExpression, CharacterExpression, request.LoggedInValidator, handler.CharacterExpressionHandler)
	hr(handler.OpCharacterCloseRangeAttack, CharacterCloseRangeAttack, request.LoggedInValidator, handler.CharacterCloseRangeAttackHandler)
	hr(handler.OpCharacterRangedAttack, CharacterRangedAttack, request.LoggedInValidator, handler.CharacterRangedAttackHandler)
	hr(handler.OpCharacterMagicAttack, CharacterMagicAttack, request.LoggedInValidator, handler.CharacterMagicAttackHandler)
	hr(handler.OpCharacterDistributeAp, CharacterDistributeAP, request.LoggedInValidator, handler.DistributeApHandler)
	hr(handler.OpCharacterDistributeSp, CharacterDistributeSP, request.LoggedInValidator, handler.DistributeSpHandler)
	hr(handler.OpCharacterHealOverTime, CharacterHealOverTime, request.LoggedInValidator, handler.HealOverTimeHandler)
	hr(handler.OpCharacterItemPickUp, CharacterItemPickUp, request.LoggedInValidator, handler.ItemPickUpHandler)
	hr(handler.OpNpcAction, NPCAction, request.LoggedInValidator, handler.HandleNPCAction)
	hr(handler.OpNpcTalkMore, NPCTalkMore, request.LoggedInValidator, handler.HandleNPCTalkMoreRequest)
	hr(handler.OpNpcTalk, NPCTalk, handler.CharacterAliveValidator, handler.HandleNPCTalkRequest)
	hr(handler.OpCharacterDamage, CharacterDamage, request.LoggedInValidator, handler.HandleCharacterDamageRequest)
	hr(handler.OpMoveItem, MoveItem, request.LoggedInValidator, handler.MoveItemHandler)
	hr(handler.OpCodeSpecialMove, SpecialMove, request.LoggedInValidator, handler.HandleSpecialMove)
	hr(handler.OpQuestAction, QuestAction, request.LoggedInValidator, handler.HandleQuestAction)
	hr(handler.OpInnerPortal, InnerPortal, request.LoggedInValidator, request.NoOpHandler)
	hr(handler.OpChangeKeyMap, ChangeKeyMap, request.LoggedInValidator, handler.ChangeKeyMapHandler)
	hr(handler.OpReactorHit, ReactorHit, request.LoggedInValidator, handler.HandleReactorHit)

	return func() map[uint16]request2.Handler {
		return handlers
	}
}
