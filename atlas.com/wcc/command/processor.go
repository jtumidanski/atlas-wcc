package command

import (
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type Producer func(s session.Model, m string) (Executor, bool)

type Executor func(l logrus.FieldLogger, span opentracing.Span) error
