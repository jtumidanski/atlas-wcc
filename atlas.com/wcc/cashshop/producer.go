package cashshop

import (
	"atlas-wcc/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type enterCashShopCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func emitEnterCashShop(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ENTER_CASH_SHOP_COMMAND")
	return func(worldId byte, channelId byte, characterId uint32) {
		e := &enterCashShopCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

type modifyWishlistCommand struct {
	CharacterId   uint32   `json:"character_id"`
	SerialNumbers []uint32 `json:"serial_numbers"`
}

func emitModifyWishlistCommand(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, serialNumbers []uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MODIFY_WISHLIST_COMMAND")
	return func(characterId uint32, serialNumbers []uint32) {
		e := &modifyWishlistCommand{
			CharacterId:   characterId,
			SerialNumbers: serialNumbers,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
