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

func requestById(id uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(partyResource, id))
}

func requestByMemberId(memberId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(characterPartyResource, memberId))
}

func createParty(worldId byte, channelId byte, characterId uint32) requests.PostRequest[attributes] {
	i := &inputDataContainer{Data: inputDataBody{
		Id:   "",
		Type: "party",
		Attributes: inputAttributes{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
		},
	}}
	return requests.MakePostRequest[attributes](partiesResource, i)
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
