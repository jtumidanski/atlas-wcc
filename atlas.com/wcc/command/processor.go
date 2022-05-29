package command

import (
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type SyntaxValidator func(c session.Model, m string) bool

type Executor func(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, m string) error
