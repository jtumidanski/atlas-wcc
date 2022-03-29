package skill

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		r, err := requestForCharacter(characterId)(l, span)
		if err != nil {
			return nil, err
		}

		ss := make([]Model, 0)
		for _, s := range r.DataList() {
			sid, err := strconv.ParseUint(s.Id, 10, 32)
			if err != nil {
				break
			}
			attr := s.Attributes
			sr := NewSkill(uint32(sid), attr.Level, attr.MasterLevel, attr.Expiration, false, false)
			ss = append(ss, sr)
		}
		return ss, nil
	}
}
