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
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			wg.Add(1)
			defer wg.Done()
			err := socket.Run(l, handlerProducer(l),
				socket.SetPort(port),
				socket.SetSessionCreator(session.Create(l, session.GetRegistry())(worldId, channelId)),
				socket.SetSessionMessageDecryptor(session.Decrypt(l, session.GetRegistry())),
				socket.SetSessionDestroyer(session.DestroyById(l, session.GetRegistry())),
			)
			if err != nil {
				l.WithError(err).Errorf("Socket service encountered error")
			}
		}()

		<-ctx.Done()
		l.Infof("Shutting down server on port 8484")
	}
}

func handlerProducer(l logrus.FieldLogger) socket.MessageHandlerProducer {
	handlers := make(map[uint16]request2.Handler)
	hr := func(op uint16, v request.MessageValidator, h request.MessageHandler) {
		handlers[op] = request.AdaptHandler(l, v, h)
	}

	hr(handler.OpCodePong, request.NoOpValidator, handler.PongHandler())
	hr(handler.OpCharacterLoggedIn, request.NoOpValidator, handler.CharacterLoggedInHandler())
	hr(handler.OpChangeMapSpecial, request.LoggedInValidator, handler.ChangeMapSpecialHandler())
	hr(handler.OpMoveCharacter, request.LoggedInValidator, handler.MoveCharacterHandler())
	hr(handler.OpChangeMap, request.LoggedInValidator, handler.ChangeMapHandler())
	hr(handler.OpMoveLife, request.LoggedInValidator, handler.MoveLifeHandler())
	hr(handler.OpGeneralChat, request.LoggedInValidator, handler.GeneralChatHandler())
	hr(handler.OpChangeChannel, request.LoggedInValidator, handler.ChangeChannelHandler())
	hr(handler.OpCharacterExpression, request.LoggedInValidator, handler.CharacterExpressionHandler())
	hr(handler.OpCharacterCloseRangeAttack, request.LoggedInValidator, handler.CharacterCloseRangeAttackHandler())
	hr(handler.OpCharacterRangedAttack, request.LoggedInValidator, handler.CharacterRangedAttackHandler())
	hr(handler.OpCharacterMagicAttack, request.LoggedInValidator, handler.CharacterMagicAttackHandler())
	hr(handler.OpCharacterDistributeAp, request.LoggedInValidator, handler.DistributeApHandler())
	hr(handler.OpCharacterDistributeSp, request.LoggedInValidator, handler.DistributeSpHandler())
	hr(handler.OpCharacterHealOverTime, request.LoggedInValidator, handler.HealOverTimeHandler())
	hr(handler.OpCharacterItemPickUp, request.LoggedInValidator, handler.ItemPickUpHandler())
	hr(handler.OpNpcAction, request.LoggedInValidator, handler.HandleNPCAction())
	hr(handler.OpNpcTalkMore, request.LoggedInValidator, handler.HandleNPCTalkMoreRequest())
	hr(handler.OpNpcTalk, handler.CharacterAliveValidator(), handler.HandleNPCTalkRequest())
	hr(handler.OpCharacterDamage, request.LoggedInValidator, handler.HandleCharacterDamageRequest())
	hr(handler.OpMoveItem, request.LoggedInValidator, handler.MoveItemHandler())
	hr(handler.OpCodeSpecialMove, request.LoggedInValidator, handler.HandleSpecialMove())
	hr(handler.OpQuestAction, request.LoggedInValidator, handler.HandleQuestAction())
	hr(handler.OpInnerPortal, request.LoggedInValidator, handler.HandleInnerPortal())
	hr(handler.OpChangeKeyMap, request.LoggedInValidator, handler.ChangeKeyMapHandler())

	return func() map[uint16]request2.Handler {
		return handlers
	}
}
