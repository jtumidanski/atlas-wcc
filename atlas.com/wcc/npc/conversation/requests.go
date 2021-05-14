package conversation

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	npcConversationServicePrefix  string = "/ms/ncs/"
	npcConversationService               = requests.BaseRequest + npcConversationServicePrefix
	npcConversationResource              = npcConversationService + "script/%d"
	characterConversationResource        = npcConversationService + "conversation/%d"
)

func HasScript(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		r, err := http.Get(fmt.Sprintf(npcConversationResource, npcId))
		if err != nil {
			l.WithError(err).Errorf("Unable to identify if npc %d has a conversation script. Assuming not.", npcId)
			return false
		}
		return r.StatusCode == http.StatusOK
	}
}

func InConversation(l logrus.FieldLogger) func(characterId uint32) bool {
	return func(characterId uint32) bool {
		r, err := http.Get(fmt.Sprintf(characterConversationResource, characterId))
		if err != nil {
			l.WithError(err).Errorf("Unable to identify if character %d is in a conversation. Assuming not.", characterId)
			return false
		}
		return r.StatusCode == http.StatusOK
	}
}
