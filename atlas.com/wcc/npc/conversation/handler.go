package conversation

import (
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/npc"
	"atlas-wcc/npc/shop"
	"atlas-wcc/session"
	"atlas-wcc/socket/request/handler"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpNpcTalk uint16 = 0x3A
const OpNpcTalkMore uint16 = 0x3C
const NPCTalk = "npc_talk"
const NPCTalkMore = "npc_talk_more"

func NPCTalkRequestHandlerProducer(l logrus.FieldLogger, worldId byte, channelId byte) handler.Producer {
	return func() (uint16, request.Handler) {
		return OpNpcTalk, handler.SpanHandlerDecorator(l, NPCTalk, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return handler.ValidatorHandler(character.CharacterAliveValidator(l, span), npcTalkRequestHandler(l, span, worldId, channelId))
		})
	}
}

func npcTalkRequestHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readNPCTalkRequest(r)

		ca, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d speaking to npc.", s.CharacterId())
			return
		}

		npcs, err := npc.GetInMapByObjectId(l, span)(ca.MapId(), p.ObjectId())
		if err != nil || len(npcs) != 1 {
			l.WithError(err).Errorf("Unable to locate npc %d in map %d.", p.ObjectId(), ca.MapId())
			return
		}
		n := npcs[0]

		if n.Id() == 9010009 {
			handleDuey(s)
			return
		}

		if n.Id() >= 9100100 && n.Id() <= 9100200 {
			handleGachapon(s)
			return
		}

		if HasScript(l)(n.Id()) {
			startConversation(l, span)(worldId, channelId, ca.MapId(), ca.Id(), n.Id(), n.ObjectId())
			return
		}
		if shop.HasShop(l, span)(n.Id()) {
			err = shop.ShowShop(l, span)(n.Id())(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to show shop for npc %d to character %d.", n.Id(), s.CharacterId())
			}
		}
	}
}

func handleGachapon(_ session.Model) {
}

func handleDuey(_ session.Model) {
}

func NPCTalkMoreRequestHandlerProducer(l logrus.FieldLogger) handler.Producer {
	return func() (uint16, request.Handler) {
		return OpNpcTalkMore, handler.SpanHandlerDecorator(l, NPCTalkMore, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return handler.ValidatorHandler(handler.LoggedInValidator(l, span), npcTalkMoreRequestHandler(l, span))
		})
	}
}

func npcTalkMoreRequestHandler(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readNPCTalkMoreRequest(r)
		if p.LastMessageType() == 2 {
			if p.Action() != 0 {
				if questInProcess(s.CharacterId()) {
					continueQuestConversation(s.CharacterId(), p)
				} else {
					setReturnText(l, span)(s.CharacterId(), p.ReturnText())
					continueConversation(l, span)(s.CharacterId(), p.Action(), p.LastMessageType(), -1)
				}
			} else if questInProcess(s.CharacterId()) {
				questDispose(s.CharacterId())
			} else {
				conversationDispose(s.CharacterId())
			}
		} else {
			if questInProcess(s.CharacterId()) {
				continueQuestConversation(s.CharacterId(), p)
			} else if InProgress(l)(s.CharacterId()) {
				continueConversation(l, span)(s.CharacterId(), p.Action(), p.LastMessageType(), p.Selection())
			}
		}
	}
}

func conversationDispose(_ uint32) {

}

func questDispose(_ uint32) {

}

func questInProcess(_ uint32) bool {
	return false
}

func continueQuestConversation(_ uint32, _ npcTalkMoreRequest) {

}
