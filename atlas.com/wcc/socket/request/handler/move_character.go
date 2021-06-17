package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpMoveCharacter uint16 = 0x29

type moveCharacterRequest struct {
	hasMovement  bool
	movementData []interface{}
	movementList []byte
}

func readMoveCharacterRequest(reader *request.RequestReader) *moveCharacterRequest {
	reader.Skip(9)

	movementDataStart := reader.Position()
	movementDataList := updatePosition(reader, 0)
	if len(movementDataList) == 0 {
		return nil
	}

	movementDataLength := reader.Position() - movementDataStart
	hasMovement := movementDataLength > 0

	movementList := make([]byte, 0)
	if hasMovement {
		reader.Seek(movementDataStart)
		for i := 0; i < movementDataLength; i++ {
			movementList = append(movementList, reader.ReadByte())
		}
	}
	return &moveCharacterRequest{hasMovement, movementDataList, movementList}
}

type absoluteMovement struct {
	X     int16
	Y     int16
	State byte
}

type relativeMovement struct {
	State byte
}

func updatePosition(reader *request.RequestReader, offset int16) []interface{} {
	commands := reader.ReadByte()

	mdl := make([]interface{}, 0)
	if commands < 1 {
		return mdl
	}

	for i := byte(0); i < commands; i++ {
		command := reader.ReadByte()
		switch command {
		case byte(0), byte(5), byte(17):
			x := reader.ReadInt16()
			y := reader.ReadInt16()
			reader.Skip(6)
			ns := reader.ReadByte()
			reader.ReadUint16()
			md := &absoluteMovement{
				X:     x,
				Y:     y + offset,
				State: ns,
			}
			mdl = append(mdl, md)
			break
		case byte(1), byte(2), byte(6), byte(12), byte(13), byte(16), byte(18), byte(19), byte(20), byte(22):
			reader.Skip(4)
			ns := reader.ReadByte()
			reader.ReadUint16()
			md := &relativeMovement{State: ns}
			mdl = append(mdl, md)
			break
		case byte(3), byte(4), byte(7), byte(8), byte(9), byte(11):
			reader.Skip(8)
			ns := reader.ReadByte()
			md := &relativeMovement{State: ns}
			mdl = append(mdl, md)
			break
		case byte(14):
			reader.Skip(9)
			break
		case byte(10):
			reader.ReadByte()
			break
		case byte(15):
			reader.Skip(12)
			ns := reader.ReadByte()
			reader.ReadUint16()
			md := &relativeMovement{State: ns}
			mdl = append(mdl, md)
			break
		case byte(21):
			reader.Skip(3)
			break
		}
	}
	return mdl
}

func MoveCharacterHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readMoveCharacterRequest(r)
		if p == nil {
			return
		}

		summary := processMovementList(p.movementData)
		producers.MoveCharacter(l)(s.WorldId(), s.ChannelId(), s.CharacterId(), summary.X, summary.Y, summary.State, p.movementList)
	}
}

func processMovementList(movementList []interface{}) movementSummary {
	ms := &movementSummary{0, 0, 0}
	for _, m := range movementList {
		ms = appendMovement(ms, m)
	}
	return *ms
}

func appendMovement(ms *movementSummary, m interface{}) *movementSummary {
	switch m.(type) {
	case *absoluteMovement:
		am := m.(*absoluteMovement)
		return &movementSummary{am.X, am.Y, am.State}
	case *relativeMovement:
		rm := m.(*relativeMovement)
		return ms.SetState(rm.State)
	}
	return ms
}

type movementSummary struct {
	X     int16
	Y     int16
	State byte
}

func (m *movementSummary) SetState(state byte) *movementSummary {
	return &movementSummary{m.X, m.Y, state}
}
