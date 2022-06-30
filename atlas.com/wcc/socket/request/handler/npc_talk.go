package handler

import (
	"atlas-wcc/account"
	"atlas-wcc/character/properties"
	npc2 "atlas-wcc/npc"
	"atlas-wcc/npc/conversation"
	"atlas-wcc/npc/shop"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpNpcTalk uint16 = 0x3A

type npcTalkRequest struct {
	objectId uint32
}

func (r npcTalkRequest) ObjectId() uint32 {
	return r.objectId
}

func readNPCTalkRequest(reader *request.RequestReader) npcTalkRequest {
	return npcTalkRequest{reader.ReadUint32()}
}

func CharacterAliveValidator(l logrus.FieldLogger, span opentracing.Span) func(s session.Model) bool {
	return func(s session.Model) bool {
		v := account.IsLoggedIn(l, span)(s.AccountId())
		if !v {
			l.Errorf("Attempting to process a [HandleNPCTalkRequest] when the account %d is not logged in.", s.SessionId())
			err := session.Announce(s, properties.WriteEnableActions(l))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return false
		}

		ca, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d speaking to npc.", s.CharacterId())
			err = session.Announce(s, properties.WriteEnableActions(l))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return false
		}

		if ca.Hp() > 0 {
			return true
		} else {
			err = session.Announce(s, properties.WriteEnableActions(l))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return false
		}
	}
}

func HandleNPCTalkRequest(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readNPCTalkRequest(r)

		ca, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d speaking to npc.", s.CharacterId())
			return
		}

		npcs, err := npc2.GetInMapByObjectId(l, span)(ca.MapId(), p.ObjectId())
		if err != nil || len(npcs) != 1 {
			l.WithError(err).Errorf("Unable to locate npc %d in map %d.", p.ObjectId(), ca.MapId())
			return
		}
		npc := npcs[0]

		if npc.Id() == 9010009 {
			handleDuey(s)
			return
		}

		if npc.Id() >= 9100100 && npc.Id() <= 9100200 {
			handleGachapon(s)
			return
		}

		if conversation.HasScript(l)(npc.Id()) {
			npc2.StartConversation(l, span)(worldId, channelId, ca.MapId(), ca.Id(), npc.Id(), npc.ObjectId())
			return
		}
		if shop.HasShop(l, span)(npc.Id()) {
			ns, err := shop.GetByNpcId(l, span)(npc.Id())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve shop for npc %d.", npc.Id())
				return
			}
			err = session.Announce(s, npc2.WriteGetNPCShop(l)(ns))
			if err != nil {
				l.WithError(err).Errorf("Unable to write shop for npc %d to character %d.", npc.Id(), s.CharacterId())
			}
		}
	}
}

func handleGachapon(_ session.Model) {

}

func handleDuey(_ session.Model) {

}
