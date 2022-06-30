package movement

import "github.com/jtumidanski/atlas-socket/request"

func readNPCAction(reader *request.RequestReader) interface{} {
	length := len(reader.GetRestAsBytes())
	if length == 6 {
		first := reader.ReadUint32()
		second := reader.ReadByte()
		third := reader.ReadByte()
		return &npcAnimationRequest{first, second, third}
	} else if length > 6 {
		bytes := reader.ReadBytes(length - 9)
		return &npcMoveRequest{bytes}
	}
	return nil
}
