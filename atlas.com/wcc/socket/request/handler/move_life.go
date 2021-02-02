package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
)

const OpMoveLife uint16 = 0xBC

type MoveLifeRequest struct {
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

func (r MoveLifeRequest) ObjectId() uint32 {
	return r.objectId
}

func (r MoveLifeRequest) RawActivity() int8 {
	return r.rawActivity
}

func (r MoveLifeRequest) PNibbles() byte {
	return r.pNibbles
}

func (r MoveLifeRequest) SkillId() uint32 {
	return r.skillId
}

func (r MoveLifeRequest) SkillLevel() uint32 {
	return r.skillLevel
}

func (r MoveLifeRequest) POption() uint16 {
	return r.pOption
}

func (r MoveLifeRequest) StartX() int16 {
	return r.startX
}

func (r MoveLifeRequest) StartY() int16 {
	return r.startY
}

func (r MoveLifeRequest) MovementData() []interface{} {
	return r.movementDataList
}

func (r MoveLifeRequest) MoveId() uint16 {
	return r.moveId
}

func (r MoveLifeRequest) MovementList() []byte {
	return r.movementList
}

func ReadMoveLifeRequest(reader *request.RequestReader) *MoveLifeRequest {
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

	return &MoveLifeRequest{objectId, moveId, pNibbles, rawActivity, skillId, skillLevel, pOption, startX, startY, hasMovement, movementDataList, movementList}
}

type MoveLifeHandler struct {
}

func (h *MoveLifeHandler) IsValid(l *log.Logger, ms *mapleSession.MapleSession) bool {
	v := processors.IsLoggedIn((*ms).AccountId())
	if !v {
		l.Printf("[ERROR] attempting to process a [MoveLifeRequest] when the account %d is not logged in.", (*ms).SessionId())
	}
	return v
}

func (h *MoveLifeHandler) HandleRequest(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
	p := ReadMoveLifeRequest(r)
	if p == nil {
		return
	}

	_, err := processors.GetMonster(p.ObjectId())
	if err != nil {
		l.Printf("[ERROR] received move life request for unknown monster %d", p.ObjectId())
		return
	}

	ra := p.RawActivity()
	pOption := p.POption()
	if ra >= 0 {
		ra = int8(int16(ra) & 0xFF >> 1)
	}

	is := h.inRangeInclusive(ra, 42, 59)

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
		producers.NewMonsterMovement(l, context.Background()).EmitMovement(p.ObjectId(), (*s).CharacterId(), nextMovementCouldBeSkill, ra, usi, usl, pOption, startX, startY, summary.X, summary.Y, summary.State, p.MovementList())
	}
}

func (h *MoveLifeHandler) inRangeInclusive(ra int8, i int8, i2 int8) bool {
	return !(ra < i) || (ra > i2)
}
