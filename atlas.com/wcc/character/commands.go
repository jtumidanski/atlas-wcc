package character

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/command"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func AwardMesoCommandProducer() command.Producer {
	return func(s session.Model, m string) (command.Executor, bool) {
		if !s.GM() {
			return nil, false
		}

		if !strings.HasPrefix(m, "@award meso") {
			return nil, false
		}
		re := regexp.MustCompile("@award meso ([^ ]*) (\\d*)")
		match := re.FindStringSubmatch(m)

		if len(match) != 3 {
			return nil, false
		}

		if len(match[1]) < 4 {
			return nil, false
		}

		amount, err := strconv.ParseUint(match[2], 10, 32)
		if err != nil {
			return nil, false
		}
		return AwardMesoCommandExecutor(match[1], int32(amount)), true
	}
}

func AwardMesoCommandExecutor(name string, amount int32) command.Executor {
	return func(l logrus.FieldLogger, span opentracing.Span) error {
		c, err := properties.GetByName(l, span)(name)
		if err != nil {
			return err
		}

		GainMeso(l, span)(c.Id(), amount)
		return nil
	}
}
