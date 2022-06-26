package cashshop

import (
	"atlas-wcc/account"
	cc "atlas-wcc/cashshop/character"
	"atlas-wcc/cashshop/character/wishlist"
	"atlas-wcc/character"
	"atlas-wcc/kafka"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameEnterCashShop = "enter_cash_shop_event"
	topicTokenEnterCashShop   = "TOPIC_ENTER_CASH_SHOP_EVENT"
)

func EnterCashShopEventConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[enterCashShopEvent](consumerNameEnterCashShop, topicTokenEnterCashShop, groupId, handleEnterCashShopEvent(wid, cid))
	}
}

type enterCashShopEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func handleEnterCashShopEvent(worldId byte, channelId byte) kafka.HandlerFunc[enterCashShopEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event enterCashShopEvent) {
		if worldId != event.WorldId || channelId != event.ChannelId {
			return
		}

		s, err := session.GetByCharacterId(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate session for character %d entering cash shop.", event.CharacterId)
			return
		}

		a, err := account.GetById(l, span)(s.AccountId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve account for character %d.", event.CharacterId)
			return
		}

		c, err := character.GetCharacterById(l, span)(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character information for %d.", event.CharacterId)
			return
		}

		scis := make([]SpecialCashItem, 0)

		err = session.Announce(WriteOpenCashShop(l)(a, c, scis))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write cash shop opening packet to character %d.", event.CharacterId)
			return
		}

		err = session.Announce(WriteCashInventory(l)(a, c))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write cash inventory to character %d.", event.CharacterId)
			return
		}

		err = session.Announce(WriteCashGifts(l)())(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write cash gifts to character %d.", event.CharacterId)
			return
		}

		wl, err := wishlist.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve wishlist for character %d.", event.CharacterId)
		}
		err = session.Announce(wishlist.WriteWishList(l)(wl, false))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write wish list to character %d.", event.CharacterId)
			return
		}

		cp, err := cc.GetById(l, span)(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve cash shop points for character %d.", event.CharacterId)
			return
		}

		_ = writeCashAmounts(l)(cp.Credit(), cp.Points(), cp.Prepaid())(s)
	}
}
