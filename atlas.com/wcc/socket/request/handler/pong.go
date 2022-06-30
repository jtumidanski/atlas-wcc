package handler

import (
	"github.com/jtumidanski/atlas-socket/request"
)

const OpCodePong uint16 = 0x18

func PongHandlerProducer() Producer {
	return func() (uint16, request.Handler) {
		return OpCodePong, ValidatorHandler(NoOpValidator, NoOpHandler)
	}
}
