package skill

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		r, err := requestForCharacter(l, span)(characterId)
		if err != nil {
			return nil, err
		}

		ss := make([]Model, 0)
		for _, s := range r.DataList() {
			sid, err := strconv.ParseUint(s.Id, 10, 32)
			if err != nil {
				break
			}
			sr := NewSkill(uint32(sid), s.Attributes.Level, s.Attributes.MasterLevel, s.Attributes.Expiration, false, false)
			ss = append(ss, sr)
		}
		return ss, nil
	}
}
