package keymap

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	keymapServicePrefix string = "/ms/cks/"
	keymapService              = requests.BaseRequest + keymapServicePrefix
	keymapResource             = keymapService + "characters/%d/keymap"
)

func getKeyMap(characterId uint32) (*DataListContainer, error) {
	ar := &DataListContainer{}
	err := requests.Get(fmt.Sprintf(keymapResource, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
