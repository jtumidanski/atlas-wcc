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

func WarpMapCommandProducer() command.Producer {
	return func(s session.Model, m string) (command.Executor, bool) {
		if !s.GM() {
			return nil, false
		}

		if !strings.HasPrefix(m, "@warp map") {
			return nil, false
		}
		re := regexp.MustCompile("@warp map (\\d*)")
		match := re.FindStringSubmatch(m)
		if len(match) != 2 {
			return nil, false
		}

		mapId, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			return nil, false
		}
		return WarpMapCommandExecutor(s.WorldId(), s.ChannelId(), s.CharacterId(), uint32(mapId)), true
	}
}

func WarpMapCommandExecutor(worldId byte, channelId byte, characterId uint32, mapId uint32) command.Executor {
	return func(l logrus.FieldLogger, span opentracing.Span) error {
		WarpRandom(l, span)(worldId, channelId, characterId, mapId)
		return nil
	}
}
