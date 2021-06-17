package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpMoveLife uint16 = 0xBC

type moveLifeRequest struct {
	objectId         uint32
	moveId           uint16
	pNibbles         byte
	rawActivity      int8
	skillId          uint32
	skillLevel       uint32
	pOption          uint16
	startX           int16
	startY           int16
	hasMovement      bool
	movementDataList []interface{}
	movementList     []byte
}

func (r moveLifeRequest) ObjectId() uint32 {
	return r.objectId
}

func (r moveLifeRequest) RawActivity() int8 {
	return r.rawActivity
}

func (r moveLifeRequest) PNibbles() byte {
	return r.pNibbles
}

func (r moveLifeRequest) SkillId() uint32 {
	return r.skillId
}

func (r moveLifeRequest) SkillLevel() uint32 {
	return r.skillLevel
}

func (r moveLifeRequest) POption() uint16 {
	return r.pOption
}

func (r moveLifeRequest) StartX() int16 {
	return r.startX
}

func (r moveLifeRequest) StartY() int16 {
	return r.startY
}

func (r moveLifeRequest) MovementData() []interface{} {
	return r.movementDataList
}

func (r moveLifeRequest) MoveId() uint16 {
	return r.moveId
}

func (r moveLifeRequest) MovementList() []byte {
	return r.movementList
}

func readMoveLifeRequest(reader *request.RequestReader) *moveLifeRequest {
	objectId := reader.ReadUint32()
	moveId := reader.ReadUint16()
	pNibbles := reader.ReadByte()
	rawActivity := reader.ReadInt8()
	skillId := uint32(reader.ReadByte() & 0xFF)
	skillLevel := uint32(reader.ReadByte() & 0xFF)
	pOption := reader.ReadUint16()
	reader.Skip(8)
	reader.ReadByte()
	reader.ReadInt32()
	startX := reader.ReadInt16()
	startY := reader.ReadInt16()

	movementDataStart := reader.Position()
	movementDataList := updatePosition(reader, -2)
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

	return &moveLifeRequest{objectId, moveId, pNibbles, rawActivity, skillId, skillLevel, pOption, startX, startY, hasMovement, movementDataList, movementList}
}

func MoveLifeHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readMoveLifeRequest(r)
		if p == nil {
			return
		}

		_, err := processors.GetMonster(p.ObjectId())
		if err != nil {
			l.WithError(err).Errorf("Received move life request for unknown monster %d", p.ObjectId())
			return
		}

		ra := p.RawActivity()
		pOption := p.POption()
		if ra >= 0 {
			ra = int8(int16(ra) & 0xFF >> 1)
		}

		is := inRangeInclusive(ra, 42, 59)

		usi := uint32(0)
		usl := uint32(0)

		nextMovementCouldBeSkill := !(is || (p.PNibbles() != 0))
		if is {
			usi = p.SkillId()
			usl = p.SkillLevel()
		} else {
			as := int32(0)
			if as < 1 {
				ra = -1
				pOption = 0
			}
		}

		startX := p.StartX()
		startY := p.StartY() - 2

		summary := processMovementList(p.MovementData())
		(*s).Announce(writer.WriteMoveMonsterResponse(p.ObjectId(), p.MoveId(), 0, false, 0, 0))

		if p.hasMovement {
			producers.MonsterMovement(l)(p.ObjectId(), (*s).CharacterId(), nextMovementCouldBeSkill, ra, usi, usl, pOption, startX, startY, summary.X, summary.Y, summary.State, p.MovementList())
		}
	}
}

func inRangeInclusive(ra int8, i int8, i2 int8) bool {
	return !(ra < i) || (ra > i2)
}
