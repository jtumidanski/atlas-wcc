package conversation

type npcTalkRequest struct {
	objectId uint32
}

func (r npcTalkRequest) ObjectId() uint32 {
	return r.objectId
}

type npcTalkMoreRequest struct {
	lastMessageType byte
	action          byte
	returnText      string
	selection       int32
}

func (r npcTalkMoreRequest) LastMessageType() byte {
	return r.lastMessageType
}

func (r npcTalkMoreRequest) Action() byte {
	return r.action
}

func (r npcTalkMoreRequest) ReturnText() string {
	return r.returnText
}

func (r npcTalkMoreRequest) Selection() int32 {
	return r.selection
}
