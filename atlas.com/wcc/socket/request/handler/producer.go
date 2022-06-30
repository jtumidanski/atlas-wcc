package handler

import request2 "github.com/jtumidanski/atlas-socket/request"

type Producer func() (uint16, request2.Handler)
