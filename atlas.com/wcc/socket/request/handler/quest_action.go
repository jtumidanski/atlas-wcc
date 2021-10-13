package handler

import (
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpQuestAction uint16 = 0x6B

type questActionRequest struct {
	action    byte
	questId   uint16
	itemId    uint32
	npcId     uint32
	selection int16
	x         int16
	y         int16
}

func (r questActionRequest) Action() byte {
	return r.action
}

func readQuestAction(r *request.RequestReader) interface{} {
	action := r.ReadByte()
	questId := r.ReadUint16()
	itemId := uint32(0)
	npcId := uint32(0)
	selection := int16(-1)
	x := int16(-1)
	y := int16(-1)

	if action == 0 {
		r.ReadInt32()
		itemId = r.ReadUint32()
	} else if action == 1 {
		npcId = r.ReadUint32()
		if len(r.GetBuffer())-r.Position() >= 4 {
			x = r.ReadInt16()
			y = r.ReadInt16()
		}
	} else if action == 2 {
		npcId = r.ReadUint32()
		if len(r.GetBuffer())-r.Position() >= 4 {
			x = r.ReadInt16()
			y = r.ReadInt16()
		}
		if len(r.GetBuffer())-r.Position() >= 2 {
			selection = r.ReadInt16()
		}
	} else if action == 4 {
		npcId = r.ReadUint32()
		if len(r.GetBuffer())-r.Position() >= 4 {
			x = r.ReadInt16()
			y = r.ReadInt16()
		}
	} else if action == 5 {
		npcId = r.ReadUint32()
		if len(r.GetBuffer())-r.Position() >= 4 {
			x = r.ReadInt16()
			y = r.ReadInt16()
		}
	}
	return &questActionRequest{
		action:    action,
		questId:   questId,
		itemId:    itemId,
		npcId:     npcId,
		selection: selection,
		x:         x,
		y:         y,
	}
}

func HandleQuestAction(l logrus.FieldLogger, _ opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readQuestAction(r)
		if val, ok := p.(*questActionRequest); ok {
			if val.Action() == 0 {
				l.Debugf("Restore lost item action.")
			} else if val.Action() == 1 {
				l.Debugf("Start Quest.")
			} else if val.Action() == 2 {
				l.Debugf("Complete Quest.")
			} else if val.Action() == 3 {
				l.Debugf("Forfeit Quest.")
			} else if val.Action() == 4 {
				l.Debugf("Start Scripted Quest.")
			} else if val.Action() == 5 {
				l.Debugf("Start Completed Quest.")
			} else {
				l.Errorf("Unhandled quest action %d.", val.Action())
			}
		} else {
			l.Warnf("Received a unhandled [NPCActionRequest]")
		}
	}
}
