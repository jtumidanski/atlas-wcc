package member

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func InPartyModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(partyId uint32) model.SliceProvider[Model] {
	return func(partyId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestMembers(partyId), makeModel)
	}
}

func GetInParty(l logrus.FieldLogger, span opentracing.Span) func(partyId uint32) ([]Model, error) {
	return func(partyId uint32) ([]Model, error) {
		return InPartyModelListProvider(l, span)(partyId)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := body.Attributes
	m := Model{id: uint32(id), characterId: att.CharacterId, worldId: att.WorldId, channelId: att.ChannelId, online: att.Online}
	return m, nil
}
