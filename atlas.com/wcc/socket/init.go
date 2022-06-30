package socket

import (
	"atlas-wcc/npc/conversation"
	"atlas-wcc/npc/movement"
	"atlas-wcc/session"
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
				err := socket.Run(l, handlerProducer(l)(worldId, channelId),
					socket.SetPort(port),
					socket.SetSessionCreator(session.Create(l)(worldId, channelId)),
					socket.SetSessionMessageDecryptor(session.Decrypt(l)),
					socket.SetSessionDestroyer(session.DestroyByIdWithSpan(l)(worldId, channelId)),
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

func handlerProducer(l logrus.FieldLogger) func(worldId byte, channelId byte) socket.MessageHandlerProducer {
	return func(worldId byte, channelId byte) socket.MessageHandlerProducer {
		handlers := make(map[uint16]request2.Handler)
		hr := func(op uint16, h request2.Handler) {
			handlers[op] = h
		}

		hr(handler.PongHandlerProducer()())
		hr(handler.CharacterLoggedInHandlerProducer(l, worldId, channelId)())
		hr(handler.ChangeMapSpecialHandlerProducer(l, worldId, channelId)())
		hr(handler.MoveCharacterHandlerProducer(l, worldId, channelId)())
		hr(handler.ChangeMapHandlerProducer(l, worldId, channelId)())
		hr(handler.MoveLifeHandlerProducer(l)())
		hr(handler.GeneralChatHandlerProducer(l, worldId, channelId)())
		hr(handler.ChangeChannelHandlerProducer(l, worldId, channelId)())
		hr(handler.CharacterExpressionHandlerProducer(l)())
		hr(handler.CharacterCloseRangeAttackHandlerProducer(l, worldId, channelId)())
		hr(handler.CharacterRangedAttackHandlerProducer(l, worldId, channelId)())
		hr(handler.CharacterMagicAttackHandlerProducer(l, worldId, channelId)())
		hr(handler.DistributeApHandlerProducer(l)())
		hr(handler.DistributeSpHandlerProducer(l)())
		hr(handler.HealOverTimeHandlerProducer(l)())
		hr(handler.ItemPickUpHandlerProducer(l)())
		hr(movement.HandleNPCActionProducer(l)())
		hr(conversation.NPCTalkMoreRequestHandlerProducer(l)())
		hr(conversation.NPCTalkRequestHandlerProducer(l, worldId, channelId)())
		hr(handler.HandleCharacterDamageRequestProducer(l)())
		hr(handler.MoveItemHandlerProducer(l, worldId, channelId)())
		hr(handler.HandleSpecialMoveProducer(l)())
		hr(handler.HandleQuestActionProducer(l)())
		hr(handler.InnerPortalHandlerProducer(l)())
		hr(handler.ChangeKeyMapHandlerProducer(l)())
		hr(handler.HandleReactorHitProducer(l, worldId, channelId)())
		hr(handler.HandlePartyOperationProducer(l, worldId, channelId)())
		hr(handler.EnterCashShopHandlerProducer(l, worldId, channelId)())
		hr(handler.TouchingCashShopHandlerProducer(l, worldId, channelId)())
		hr(handler.CashShopOperationHandlerProducer(l)())

		return func() map[uint16]request2.Handler {
			return handlers
		}
	}
}
