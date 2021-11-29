package party

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	partyRegistryServicePrefix string = "/ms/party/"
	partyRegistryService              = requests.BaseRequest + partyRegistryServicePrefix
	partiesResource                   = partyRegistryService + "parties"
	partyResource                     = partiesResource + "/%d"
	partyMembersResource              = partyResource + "/members"
	partyMemberResource               = partyMembersResource + "/%d"
	charactersResource                = partyRegistryService + "characters"
	characterResource                 = charactersResource + "/%d"
	characterPartyResource            = characterResource + "/party"
)

type Request func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error)

func makeRequest(url string) Request {
	return func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(url, ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestById(id uint32) Request {
	return makeRequest(fmt.Sprintf(partyResource, id))
}

func requestByMemberId(memberId uint32) Request {
	return makeRequest(fmt.Sprintf(characterPartyResource, memberId))
}

func createParty(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) error {
	return func(worldId byte, channelId byte, characterId uint32) error {
		i := &inputDataContainer{Data: inputDataBody{
			Id:   "",
			Type: "party",
			Attributes: inputAttributes{
				WorldId:     worldId,
				ChannelId:   channelId,
				CharacterId: characterId,
			},
		}}

		er := &requests.ErrorListDataContainer{}
		return requests.Post(l, span)(partiesResource, i, nil, er)
	}
}

func leaveParty(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, partyId uint32, characterId uint32) error {
	return func(worldId byte, channelId byte, partyId uint32, characterId uint32) error {
		i := &inputDataContainer{Data: inputDataBody{
			Id:   "",
			Type: "party",
			Attributes: inputAttributes{
				WorldId:     worldId,
				ChannelId:   channelId,
				CharacterId: characterId,
			},
		}}

		return requests.Delete(l, span)(fmt.Sprintf(partyMemberResource, partyId, characterId), i)
	}
}
