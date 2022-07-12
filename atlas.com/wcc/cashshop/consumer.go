package cashshop

import (
	"atlas-wcc/account"
	cc "atlas-wcc/cashshop/character"
	"atlas-wcc/cashshop/character/wishlist"
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/kafka"
	"atlas-wcc/server"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameEnterCashShop          = "enter_cash_shop_event"
	consumerNameCashShopEntryRejection = "cash_shop_entry_rejection_event"
	topicTokenEnterCashShop            = "TOPIC_ENTER_CASH_SHOP_EVENT"
	topicTokenCashShopEntryRejection   = "TOPIC_CASH_SHOP_ENTRY_REJECTION_EVENT"
)

func EnterEventConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[enterEvent](consumerNameEnterCashShop, topicTokenEnterCashShop, groupId, handleEnterEvent(wid, cid))
	}
}

type enterEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func handleEnterEvent(worldId byte, channelId byte) kafka.HandlerFunc[enterEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event enterEvent) {
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

		err = session.Announce(s, WriteOpenCashShop(l)(a, c, scis))
		if err != nil {
			l.WithError(err).Errorf("Unable to write cash shop opening packet to character %d.", event.CharacterId)
			return
		}

		err = session.Announce(s, WriteCashInventory(l)(a, c))
		if err != nil {
			l.WithError(err).Errorf("Unable to write cash inventory to character %d.", event.CharacterId)
			return
		}

		err = session.Announce(s, WriteCashGifts(l)())
		if err != nil {
			l.WithError(err).Errorf("Unable to write cash gifts to character %d.", event.CharacterId)
			return
		}

		wl, err := wishlist.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve wishlist for character %d.", event.CharacterId)
		}
		err = session.Announce(s, wishlist.WriteWishList(l)(wl, false))
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

type entryRejectionEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

func EntryRejectionConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[entryRejectionEvent](consumerNameCashShopEntryRejection, topicTokenCashShopEntryRejection, groupId, handleRejectionEvent(wid, cid))
	}
}

func handleRejectionEvent(worldId byte, channelId byte) kafka.HandlerFunc[entryRejectionEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event entryRejectionEvent) {
		if worldId != event.WorldId || channelId != event.ChannelId {
			return
		}

		var op []byte
		if event.MessageType == "POP_UP" {
			op = server.WritePopup(l)(event.Message)
		} else if event.MessageType == "PINK_TEXT" {
			op = server.WritePinkText(l)(event.Message)
		} else {
			l.Errorf("Unhandled MessageType %s provided to inform character %s of cash shop entry rejection.", event.MessageType, event.CharacterId)
			return
		}

		session.IfPresentByCharacterId(event.CharacterId, session.AnnounceOperator(op, properties.WriteEnableActions(l)))
	}
}
