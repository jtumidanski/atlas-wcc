package cashshop

import (
	"atlas-wcc/cashshop/character"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func RequestCashShopEntry(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		emitEnterCashShop(l, span)(worldId, channelId, characterId)
	}
}

func UpdateCashAmounts(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		c, err := character.GetById(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve cash amounts for character %d.", characterId)
			return
		}

		session.ForSessionByCharacterId(characterId, writeCashAmounts(l)(c.Credit(), c.Points(), c.Prepaid()))
	}
}

func writeCashAmounts(l logrus.FieldLogger) func(credit uint32, points uint32, prepaid uint32) model.Operator[session.Model] {
	return func(credit uint32, points uint32, prepaid uint32) model.Operator[session.Model] {
		return func(s session.Model) error {
			err := session.Announce(WriteCashInformation(l)(credit, points, prepaid))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to write cash information to character %d.", s.CharacterId())
			}
			return err
		}
	}
}

func RequestItemPurchase(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, cashIndex uint32, serialNumber uint32) {
	return func(characterId uint32, cashIndex uint32, serialNumber uint32) {
		emitPurchaseItemCommand(l, span)(characterId, cashIndex, serialNumber)
	}
}

func ModifyWishlist(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, serialNumbers []uint32) {
	return func(characterId uint32, serialNumbers []uint32) {
		emitModifyWishlistCommand(l, span)(characterId, serialNumbers)
	}
}
