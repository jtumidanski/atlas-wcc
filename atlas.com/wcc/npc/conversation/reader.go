package conversation

import "github.com/jtumidanski/atlas-socket/request"

func readNPCTalkRequest(reader *request.RequestReader) npcTalkRequest {
	return npcTalkRequest{reader.ReadUint32()}
}

func readNPCTalkMoreRequest(reader *request.RequestReader) npcTalkMoreRequest {
	lastMessageType := reader.ReadByte()
	action := reader.ReadByte()
	returnText := ""
	selection := int32(-1)

	if lastMessageType == 2 {
		if action != 0 {
			returnText = reader.ReadAsciiString()
		}
	} else {
		if len(reader.GetRestAsBytes()) >= 4 {
			selection = reader.ReadInt32()
		} else {
			selection = int32(reader.ReadByte())
		}
	}
	return npcTalkMoreRequest{lastMessageType, action, returnText, selection}
}
