package _map

import (
	"atlas-wcc/command"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func WarpMapCommandSyntaxValidator() command.SyntaxValidator {
	return func(s session.Model, m string) bool {
		if !s.GM() {
			return false
		}

		if !strings.HasPrefix(m, "@warp map") {
			return false
		}
		re := regexp.MustCompile("@warp map (\\d*)")
		match := re.FindStringSubmatch(m)
		_, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			return false
		}
		return true
	}
}

func WarpMapCommandExecutor() command.Executor {
	return func(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, m string) error {
		return func(s session.Model, m string) error {
			re := regexp.MustCompile("@warp map (\\d*)")
			match := re.FindStringSubmatch(m)
			mapId, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return err
			}

			WarpRandom(l, span)(s.WorldId(), s.ChannelId(), s.CharacterId(), uint32(mapId))
			return nil
		}
	}
}
