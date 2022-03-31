package party

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ByMemberIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(memberId uint32) model.Provider[Model] {
	return func(memberId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestByMemberId(memberId), makeModel)
	}
}

func GetByMemberId(l logrus.FieldLogger, span opentracing.Span) func(memberId uint32) (Model, error) {
	return func(memberId uint32) (Model, error) {
		return ByMemberIdModelProvider(l, span)(memberId)()
	}
}

func Create(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		_, _, err := createParty(worldId, channelId, characterId)(l, span)
		if err != nil {
			l.WithError(err).Errorf("Unable to create party for character %d.", characterId)
			return
		}
	}
}

func Leave(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		p, err := GetByMemberId(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate party for character %d.", characterId)
			return
		}

		err = leaveParty(l, span)(worldId, channelId, p.Id(), characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to create party for character %d.", characterId)
			return
		}
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := body.Attributes
	m := Model{id: uint32(id), leaderId: att.LeaderId}
	return m, nil
}
