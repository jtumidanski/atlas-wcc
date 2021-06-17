package handler

import (
	"atlas-wcc/session"
	"atlas-wcc/socket/request"
	request2 "github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpInnerPortal uint16 = 0x65


func HandleInnerPortal() request.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request2.RequestReader) {
	}
}