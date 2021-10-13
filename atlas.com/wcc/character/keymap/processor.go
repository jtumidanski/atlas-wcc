package keymap

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetKeyMap(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]*Model, error) {
	return func(characterId uint32) ([]*Model, error) {
		r, err := requestKeyMap(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve keymap for character.")
			return nil, err
		}

		keys := make([]*Model, 0)
		for _, data := range r.Data {
			k, err := makeKey(data)
			if err != nil {
				l.WithError(err).Errorf("Unable to create keybinding for key %s.", data.Id)
				return nil, err
			}
			keys = append(keys, k)
		}
		return keys, nil
	}
}

func makeKey(k DataBody) (*Model, error) {
	id, err := strconv.Atoi(k.Id)
	if err != nil {
		return nil, err
	}

	return &Model{
		id:      uint32(id),
		key:     k.Attributes.Key,
		theType: k.Attributes.Type,
		action:  k.Attributes.Action,
	}, nil
}
