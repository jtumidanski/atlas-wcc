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

func requestKeyMap(characterId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(keymapResource, characterId))
}
