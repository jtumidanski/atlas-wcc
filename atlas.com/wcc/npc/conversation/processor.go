package conversation

import "github.com/sirupsen/logrus"

func HasScript(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		return hasScript(l)(npcId)
	}
}

func InProgress(l logrus.FieldLogger) func(characterId uint32) bool {
	return func(characterId uint32) bool {
		return inConversation(l)(characterId)
	}
}