package keymap

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	keymapServicePrefix string = "/ms/cks/"
	keymapService              = requests.BaseRequest + keymapServicePrefix
	keymapResource             = keymapService + "characters/%d/keymap"
)

func requestKeyMap(l logrus.FieldLogger) func(characterId uint32) (*DataListContainer, error) {
	return func(characterId uint32) (*DataListContainer, error) {
		ar := &DataListContainer{}
		err := requests.Get(l)(fmt.Sprintf(keymapResource, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
